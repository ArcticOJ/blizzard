package judge

import (
	"blizzard/cache/stores"
	"blizzard/config"
	"blizzard/db"
	"blizzard/db/models/contest"
	"blizzard/logger"
	"blizzard/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	amqp "github.com/rabbitmq/amqp091-go"
	amqp2 "github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"github.com/uptrace/bun"
	"github.com/vmihailenco/msgpack"
	"golang.org/x/sync/semaphore"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var ResponseWorker *Worker

type (
	Worker struct {
		c   *http.Client
		s   *semaphore.Weighted
		ctx context.Context
		// subscribers
		sm map[uint32][]chan<- interface{}
		// sub channel
		sc chan Subscription
		// unsub channel
		usc chan Subscription
		// pub channel
		pc chan result
		// destroy channel
		dc chan uint32
		// result message queue
		errChan    <-chan *amqp.Error
		returnChan <-chan amqp.Return
		rmq        amqp.Queue
		mqChan     *amqp.Channel
		mqConn     *amqp.Connection
		env        *stream.Environment
	}

	Subscription struct {
		id uint32
		c  chan<- interface{}
	}

	result struct {
		id   uint32
		ttl  uint16
		data interface{}
	}

	rmqApiResponse struct {
		Destination string                 `json:"destination"`
		Arguments   map[string]interface{} `json:"arguments"`
		RoutingKey  string                 `json:"routing_key"`
	}

	res struct {
		Headers map[string]interface{}
		Body    interface{}
	}
)

func NewObserver(ctx context.Context) *Worker {
	return &Worker{
		c: &http.Client{
			Timeout: time.Second,
		},
		s:          semaphore.NewWeighted(8),
		ctx:        ctx,
		sm:         make(map[uint32][]chan<- interface{}),
		sc:         make(chan Subscription, 1),
		pc:         make(chan result, 1),
		usc:        make(chan Subscription, 1),
		dc:         make(chan uint32, 1),
		errChan:    make(<-chan *amqp.Error, 1),
		returnChan: make(<-chan amqp.Return, 1),
	}
}

func (w *Worker) Connect() {
	var e error
	if w.mqConn != nil {
		w.mqConn.Close()
		w.mqChan.Close()
	}
	conf := config.Config.RabbitMQ
	w.mqConn, e = amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s", conf.Username, conf.Password, net.JoinHostPort(conf.Host, fmt.Sprint(conf.Port))), amqp.Config{
		Heartbeat: time.Second,
		Vhost:     conf.VHost,
	})
	logger.Panic(e, "failed to establish a connection to rabbitmq")
	w.mqChan, e = w.mqConn.Channel()
	logger.Panic(e, "failed to open a channel for queue")
	logger.Panic(w.mqChan.ExchangeDeclare("submissions", "direct", true, false, false, false, amqp.Table{
		"x-consumer-timeout": 3600000,
	}), "failed to declare exchange for submissions")
	logger.Panic(w.mqChan.Qos(100, 0, false), "failed to set qos")
	logger.Panic(w.mqChan.ExchangeDeclare("results", "fanout", true, false, false, false, nil), "could not declare exchange for results")
	w.errChan = w.mqConn.NotifyClose(make(chan *amqp.Error, 1))
	w.returnChan = w.mqChan.NotifyReturn(make(chan amqp.Return, 1))
	w.RecoverResults()
}

func (w *Worker) RecoverResults() {
	conf := config.Config.RabbitMQ
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/api/exchanges/%s/results/bindings/source", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost)), nil)
	req.SetBasicAuth(conf.Username, conf.Password)
	if res, e := w.c.Do(req); e == nil {
		var r []rmqApiResponse
		if json.NewDecoder(res.Body).Decode(&r) != nil {
			return
		}
		for _, _r := range r {
			_id, ok := _r.Arguments["x-id"]
			if !ok || !stores.Pending.IsPending(w.ctx, uint32(_id.(float64))) {
				w.mqChan.QueueDelete(_r.Destination, true, false, true)
			} else {
				if _e := w.Consume(uint32(_id.(float64)), _r.Destination); e != nil {
					logger.Blizzard.Error().Err(_e).Msgf("could not recover results for submission '%v'", _id)
					continue
				}
			}
		}
	} else {
		logger.Blizzard.Error().Err(e).Msg("failed to get available streams for results")
	}
}

// TODO: implement auto reconnect

func (w *Worker) CreateStream() {
	var e error
	conf := config.Config.RabbitMQ
	w.env, e = stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost(conf.Host).
			SetPort(int(conf.StreamPort)).
			SetUser(conf.Username).
			SetPassword(conf.Password).
			SetVHost(conf.VHost))
	if e != nil {
		logger.Panic(e, "failed to start a stream connection")
	}
}

func (w *Worker) commitToDb(id uint32, res *contest.FinalResult) {
	db.Database.RunInTx(w.ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, e := tx.NewUpdate().Model((*contest.Submission)(nil)).Where("id = ?", id).Set("result = ?", res).Returning("NULL").Exec(w.ctx); e != nil {
			return e
		}
		return nil
	})
}

func (w *Worker) Enqueue(sub *Submission, t time.Time) error {
	var e error
	b, e := msgpack.Marshal(sub)
	if e != nil {
		return e
	}
	timestamp := time.Now().UTC().UnixMilli()
	name := fmt.Sprintf("%d_%d", sub.ID, timestamp)
	if e = w.env.DeclareStream(name, &stream.StreamOptions{
		MaxAge: time.Hour * 24,
	}); e != nil {
		return e
	}
	e = w.mqChan.QueueBind(name, "", "results", false, amqp.Table{
		"x-id":        int64(sub.ID),
		"x-timestamp": timestamp,
		"x-capacity":  int(sub.TestCount + 1),
	})
	if e != nil {
		w.env.DeleteStream(name)
		return e
	}
	e = w.Consume(sub.ID, name)
	if e != nil {
		w.env.DeleteStream(name)
		return e
	}
	conf, e := w.mqChan.PublishWithDeferredConfirmWithContext(
		w.ctx,
		"submissions",
		sub.Language,
		true,
		false,
		amqp.Publishing{
			Timestamp:     t,
			ReplyTo:       name,
			DeliveryMode:  amqp.Persistent,
			CorrelationId: strconv.FormatUint(uint64(sub.ID), 10),
			ContentType:   "application/msgpack",
			Body:          b,
		})
	if e != nil {
		w.env.DeleteStream(name)
	}
	return stores.Pending.Set(w.ctx, sub.ID, conf.DeliveryTag, sub.TestCount, uint16(math.Ceil(float64(sub.Constraints.TimeLimit))))

}

func (w *Worker) Cancel(ctx context.Context, id uint32) error {
	if !stores.Pending.IsPending(ctx, id) {
		return errors.New("no submission with matching ID")
	}
	if tag, ok := stores.Pending.Get(ctx, id); ok {
		return w.mqChan.Reject(tag, false)
	}
	return errors.New("could not cancel specified submission")
}

func (w *Worker) Consume(id uint32, name string) error {
	lastNonAcVerdict := contest.None
	_, e := w.env.NewConsumer(name, func(ctx stream.ConsumerContext, msg *amqp2.Message) {
		var r res
		if msgpack.Unmarshal(msg.Data[0], &r) != nil {
			return
		}
		if r.Headers["type"] == "final" {
			ctx.Consumer.Close()
			var _r FinalResult
			if mapstructure.Decode(r.Body, &_r) != nil {
				return
			}
			_r.LastNonACVerdict = lastNonAcVerdict
			w.publish(id, _r, 0)
			w.env.DeleteStream(name)
			stores.Pending.Delete(w.ctx, id)
		} else {
			var _r CaseResult
			if mapstructure.Decode(r.Body, &_r) != nil {
				return
			}
			if _r.Verdict != Accepted {
				lastNonAcVerdict = resolveVerdict(_r.Verdict)
			}
			var ttl uint16 = 0
			if _ttl, _ok := r.Headers["ttl"].(int32); _ok && _ttl > 0 {
				ttl = uint16(_ttl)
			}
			w.publish(id, _r, ttl)
		}
	}, stream.NewConsumerOptions().SetOffset(stream.OffsetSpecification{}.First()).SetCRCCheck(false))
	return e
}

func (w *Worker) Reconnect() {
	for {
		select {
		case <-w.errChan:
			logger.Blizzard.Debug().Msg("reconnect")
			w.Connect()
		}
	}
}

func (w *Worker) Work() {
	w.CreateStream()
	w.Connect()
	go w.Reconnect()
	for {
		select {
		// on destroy
		case id := <-w.dc:
			stores.Pending.Delete(w.ctx, id)
			a, ok := w.sm[id]
			if !ok {
				continue
			}
			for i := range a {
				close(a[i])
				a[i] = nil
			}
			w.sm[id] = nil
			delete(w.sm, id)
		// on sub
		case s := <-w.sc:
			w.sm[s.id] = append(w.sm[s.id], s.c)
		// on unsub
		case u := <-w.usc:
			m, ok := w.sm[u.id]
			if !ok {
				continue
			}
			m = utils.ArrayRemove(m, func(r chan<- interface{}) bool {
				if r == u.c {
					close(u.c)
					u.c = nil
					return true
				}
				return false
			})
			w.sm[u.id] = m
		// on pub
		case msg := <-w.pc:
			a, ok := w.sm[msg.id]
			isFinal := false
			if !ok {
				continue
			}
			var d interface{} = nil
			switch r := msg.data.(type) {
			case CaseResult:
				cr := contest.CaseResult{
					Message:  r.Message,
					Verdict:  resolveVerdict(r.Verdict),
					Memory:   r.Memory,
					Duration: r.Duration,
				}
				d = cr
			case FinalResult:
				fr := resolveFinalResult(r)
				d = fr
				w.commitToDb(msg.id, fr)
				isFinal = true
			}
			if d != nil {
				for i := range a {
					select {
					case a[i] <- d:
					}
				}
			}
			if isFinal {
				w.DestroyObserver(msg.id)
			}
		}
	}
}

func (w *Worker) publish(id uint32, data interface{}, ttl uint16) {
	w.pc <- result{id: id, data: data, ttl: ttl}
}

func (w *Worker) DestroyObserver(id uint32) {
	w.dc <- id
}

func (w *Worker) Observe(id uint32, c chan<- interface{}) (s Subscription) {
	s = Subscription{
		id: id,
		c:  c,
	}
	w.sc <- s
	return
}

func (w *Worker) StopObserving(s Subscription) {
	w.usc <- s
}
