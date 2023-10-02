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

var ResponseObserver *Observer

type (
	Observer struct {
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

func NewObserver(ctx context.Context) *Observer {
	return &Observer{
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

func (o *Observer) CheckAvailability(language string, ctx context.Context) bool {
	if o.s.Acquire(ctx, 1) != nil {
		return false
	}
	defer o.s.Release(1)
	conf := config.Config.RabbitMQ
	req, e := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s:%d/api/exchanges/%s/submissions/bindings/source", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost)), nil)
	if e != nil {
		return false
	}
	req.SetBasicAuth(conf.Username, conf.Password)
	if r, e := o.c.Do(req); e == nil {
		var _r []rmqApiResponse
		if json.NewDecoder(r.Body).Decode(&_r) != nil {
			return false
		}
		return utils.ArrayFind(_r, func(response rmqApiResponse) bool {
			return response.RoutingKey == language
		})
	}
	return false
}

func (o *Observer) Connect() {
	var e error
	if o.mqConn != nil {
		o.mqConn.Close()
		o.mqChan.Close()
	}
	conf := config.Config.RabbitMQ
	o.mqConn, e = amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s", conf.Username, conf.Password, net.JoinHostPort(conf.Host, fmt.Sprint(conf.Port))), amqp.Config{
		Heartbeat: time.Second,
		Vhost:     conf.VHost,
	})
	logger.Panic(e, "failed to establish a connection to rabbitmq")
	o.mqChan, e = o.mqConn.Channel()
	logger.Panic(e, "failed to open a channel for queue")
	logger.Panic(o.mqChan.ExchangeDeclare("submissions", "direct", true, false, false, false, amqp.Table{
		"x-consumer-timeout": 3600000,
	}), "failed to declare exchange for submissions")
	logger.Panic(o.mqChan.Confirm(false), "could not put queue channel to confirm mode")
	logger.Panic(o.mqChan.Qos(100, 0, false), "failed to set qos")
	logger.Panic(o.mqChan.ExchangeDeclare("results", "fanout", true, false, false, false, nil), "could not declare exchange for results")
	o.errChan = o.mqConn.NotifyClose(make(chan *amqp.Error, 1))
	o.returnChan = o.mqChan.NotifyReturn(make(chan amqp.Return, 1))
	o.RecoverResults()
}

func (o *Observer) RecoverResults() {
	conf := config.Config.RabbitMQ
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/api/exchanges/%s/results/bindings/source", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost)), nil)
	req.SetBasicAuth(conf.Username, conf.Password)
	if res, e := o.c.Do(req); e == nil {
		var r []rmqApiResponse
		if json.NewDecoder(res.Body).Decode(&r) != nil {
			return
		}
		for _, _r := range r {
			_id, ok := _r.Arguments["x-id"]
			if !ok || !stores.Pending.IsPending(o.ctx, uint32(_id.(float64))) {
				o.mqChan.QueueDelete(_r.Destination, true, false, true)
			} else {
				if _e := o.Consume(uint32(_id.(float64)), _r.Destination); e != nil {
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

func (o *Observer) CreateStream() {
	var e error
	conf := config.Config.RabbitMQ
	o.env, e = stream.NewEnvironment(
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

func (o *Observer) commitToDb(id uint32, res *contest.FinalResult) {
	db.Database.RunInTx(o.ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, e := tx.NewUpdate().Model((*contest.Submission)(nil)).Where("id = ?", id).Set("result = ?", res).Returning("NULL").Exec(o.ctx); e != nil {
			return e
		}
		return nil
	})
}

func (o *Observer) Enqueue(sub *Submission, t time.Time) error {
	var e error
	b, e := msgpack.Marshal(sub)
	if e != nil {
		return e
	}
	timestamp := time.Now().UTC().UnixMilli()
	name := fmt.Sprintf("%d_%d", sub.ID, timestamp)
	if e = o.env.DeclareStream(name, &stream.StreamOptions{
		MaxAge: time.Hour * 24,
	}); e != nil {
		return e
	}
	e = o.mqChan.QueueBind(name, "", "results", false, amqp.Table{
		"x-id":        int64(sub.ID),
		"x-timestamp": timestamp,
		"x-capacity":  int(sub.TestCount + 1),
	})
	if e != nil {
		o.env.DeleteStream(name)
		return e
	}
	e = o.Consume(sub.ID, name)
	if e != nil {
		o.env.DeleteStream(name)
		return e
	}
	conf, e := o.mqChan.PublishWithDeferredConfirmWithContext(
		o.ctx,
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
		o.env.DeleteStream(name)
	}
	return stores.Pending.Set(o.ctx, sub.ID, conf.DeliveryTag, sub.TestCount, uint16(math.Ceil(float64(sub.Constraints.TimeLimit))))

}

func (o *Observer) Cancel(ctx context.Context, id uint32) error {
	if !stores.Pending.IsPending(ctx, id) {
		return errors.New("no submission with matching ID")
	}
	if tag, ok := stores.Pending.Get(ctx, id); ok {
		return o.mqChan.Reject(tag, false)
	}
	return errors.New("could not cancel specified submission")
}

func (o *Observer) Consume(id uint32, name string) error {
	lastNonAcVerdict := contest.None
	_, e := o.env.NewConsumer(name, func(ctx stream.ConsumerContext, msg *amqp2.Message) {
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
			o.publish(id, _r, 0)
			o.env.DeleteStream(name)
			stores.Pending.Delete(o.ctx, id)
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
			o.publish(id, _r, ttl)
		}
	}, stream.NewConsumerOptions().SetOffset(stream.OffsetSpecification{}.First()).SetCRCCheck(false))
	return e
}

func (o *Observer) Reconnect() {
	for {
		select {
		case <-o.errChan:
			logger.Blizzard.Debug().Msg("reconnect")
			o.Connect()
		}
	}
}

func (o *Observer) Work() {
	o.CreateStream()
	o.Connect()
	go o.Reconnect()
	for {
		select {
		// on destroy
		case id := <-o.dc:
			fmt.Printf("destroying %d\n", id)
			stores.Pending.Delete(o.ctx, id)
			a, ok := o.sm[id]
			if !ok {
				continue
			}
			for i := range a {
				close(a[i])
				a[i] = nil
			}
			o.sm[id] = nil
			delete(o.sm, id)
		// on sub
		case s := <-o.sc:
			o.sm[s.id] = append(o.sm[s.id], s.c)
		// on unsub
		case u := <-o.usc:
			m, ok := o.sm[u.id]
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
			o.sm[u.id] = m
		// on pub
		case msg := <-o.pc:
			a, ok := o.sm[msg.id]
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
				o.commitToDb(msg.id, fr)
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
				o.DestroyObserver(msg.id)
			}
		}
	}
}

func (o *Observer) publish(id uint32, data interface{}, ttl uint16) {
	o.pc <- result{id: id, data: data, ttl: ttl}
}

func (o *Observer) DestroyObserver(id uint32) {
	o.dc <- id
}

func (o *Observer) Observe(id uint32, c chan<- interface{}) (s Subscription) {
	s = Subscription{
		id: id,
		c:  c,
	}
	o.sc <- s
	return
}

func (o *Observer) StopObserve(s Subscription) {
	o.usc <- s
}
