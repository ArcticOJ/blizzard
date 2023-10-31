package judge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/cache/stores"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/Jeffail/tunny"
	"github.com/mitchellh/mapstructure"
	cmap "github.com/orcaman/concurrent-map/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	amqp2 "github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"github.com/uptrace/bun"
	"github.com/vmihailenco/msgpack"
	"math"
	"net"
	"net/http"
	"net/url"
	rt "runtime"
	"slices"
	"strconv"
	"time"
)

var Worker *worker

type (
	worker struct {
		c   *http.Client
		ctx context.Context
		// subscribers
		// use nested hashmap for result chan for easier removal
		sm cmap.ConcurrentMap[uint32, []chan interface{}]
		// result message queue
		errChan    <-chan *amqp.Error
		returnChan <-chan amqp.Return
		rmq        amqp.Queue
		mqChan     *amqp.Channel
		mqConn     *amqp.Connection
		env        *stream.Environment
		pool       *tunny.Pool
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

	consumeParams struct {
		id   uint32
		name string
	}
)

func Init(ctx context.Context) {
	Worker = &worker{
		c: &http.Client{
			Timeout: time.Second,
		},
		ctx: ctx,
		sm: cmap.NewWithCustomShardingFunction[uint32, []chan interface{}](func(key uint32) uint32 {
			return key
		}),
		errChan:    make(<-chan *amqp.Error, 1),
		returnChan: make(<-chan amqp.Return, 1),
	}
	Worker.pool = tunny.NewFunc(rt.NumCPU(), func(i interface{}) interface{} {
		if p, ok := i.(consumeParams); ok {
			return Worker.consume(p.id, p.name)
		}
		return nil
	})
}

func (w *worker) Connect() {
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
		"x-consumer-timeout": time.Hour.Milliseconds(),
	}), "failed to declare exchange for submissions")
	logger.Panic(w.mqChan.Qos(100, 0, false), "failed to set qos")
	logger.Panic(w.mqChan.ExchangeDeclare("results", "direct", true, false, false, false, nil), "could not declare exchange for results")
	//logger.Panic(w.mqChan.Confirm(false), "failed to put channel to confirm mode")
	w.errChan = w.mqConn.NotifyClose(make(chan *amqp.Error, 1))
	w.returnChan = w.mqChan.NotifyReturn(make(chan amqp.Return, 1))
	w.RecoverResults()
}

func (w *worker) RecoverResults() {
	conf := config.Config.RabbitMQ
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/api/exchanges/%s/results/bindings/source", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost)), nil)
	req.SetBasicAuth(conf.Username, conf.Password)
	if r, e := w.c.Do(req); e == nil {
		var mqr []rmqApiResponse
		if json.NewDecoder(r.Body).Decode(&mqr) != nil {
			return
		}
		for _, _r := range mqr {
			_id, ok := _r.Arguments["x-id"]
			if !ok || !stores.Submissions.IsPending(w.ctx, uint32(_id.(float64))) {
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

func (w *worker) CreateStream() {
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

func (w *worker) commitToDb(id uint32, cases []contest.CaseResult, fr finalResult, v contest.Verdict, p float64) {
	logger.Panic(db.Database.RunInTx(w.ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, e := tx.NewUpdate().Model(&contest.Submission{
			ID:             id,
			Results:        cases,
			Verdict:        v,
			Points:         p,
			CompilerOutput: fr.CompilerOutput,
		}).WherePK().Column("results", "verdict", "points", "compiler_output").Returning("NULL").Exec(w.ctx)
		return e
	}), "tx")
}

func (w *worker) Enqueue(sub *Submission, t time.Time) error {
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
	e = w.mqChan.PublishWithContext(
		w.ctx,
		"submissions",
		sub.Language,
		true,
		false,
		amqp.Publishing{
			Timestamp:     t,
			ReplyTo:       name,
			DeliveryMode:  amqp.Transient,
			CorrelationId: strconv.FormatUint(uint64(sub.ID), 10),
			ContentType:   "application/msgpack",
			Body:          b,
		})
	if e != nil {
		w.env.DeleteStream(name)
		return e
	}
	return stores.Submissions.SetPending(w.ctx, sub.ID, 0, sub.TestCount, uint16(math.Ceil(float64(sub.Constraints.TimeLimit))))
}

func (w *worker) Cancel(ctx context.Context, id uint32) error {
	if !stores.Submissions.IsPending(ctx, id) {
		return errors.New("no submission with matching ID")
	}
	if tag, ok := stores.Submissions.GetPendingTag(ctx, id); ok {
		return w.mqChan.Reject(tag, false)
	}
	return errors.New("could not cancel specified submission")
}

// TODO: figure out a way to avoid race condition as two judges may judge a submission concurrently, wasting resources.

func (w *worker) consume(id uint32, name string) bool {
	lastNonAcVerdict := contest.None
	_, e := w.env.NewConsumer(name, func(ctx stream.ConsumerContext, msg *amqp2.Message) {
		var r res
		if msgpack.Unmarshal(msg.Data[0], &r) != nil {
			return
		}
		switch r.Headers["type"] {
		case "final":
			ctx.Consumer.Close()
			var _r finalResult
			if mapstructure.Decode(r.Body, &_r) != nil {
				return
			}
			_r.LastNonACVerdict = lastNonAcVerdict
			w.publish(id, math.MaxUint16, _r)
			w.env.DeleteStream(name)
			break
		case "announcement":
			if cid, ok := r.Body.(uint16); ok {
				a := announcement{
					Type: "compile",
					ID:   cid,
				}
				if cid > 0 {
					a.Type = "case"
				}
				w.publish(id, cid, a)
			}
			break
		default:
			var _r caseResult
			if mapstructure.Decode(r.Body, &_r) != nil {
				return
			}
			if _r.Verdict != Accepted {
				lastNonAcVerdict = resolveVerdict(_r.Verdict)
			}
			var cid uint16 = math.MaxUint16
			if _cid, _ok := r.Headers["case-id"].(int32); _ok {
				cid = uint16(_cid)
			}
			w.publish(id, cid, _r)
			break
		}
	}, stream.NewConsumerOptions().SetOffset(stream.OffsetSpecification{}.First()).SetCRCCheck(false))
	return e == nil
}

func (w *worker) Consume(id uint32, name string) error {
	if e, ok := w.pool.Process(consumeParams{
		id:   id,
		name: name,
	}).(bool); ok && !e {
		return fmt.Errorf("could not start consuming queue '%s'", name)
	}
	return nil
}

func (w *worker) Reconnect() {
	for {
		select {
		case <-w.errChan:
			logger.Blizzard.Debug().Msg("reconnect")
			w.Connect()
		}
	}
}

func (w *worker) Work() {
	w.CreateStream()
	w.Connect()
	w.Reconnect()
}

func (w *worker) publish(id uint32, cid uint16, data interface{}) {
	subscribers, _ := w.sm.Get(id)
	var d interface{} = nil
	switch r := data.(type) {
	case caseResult:
		cr := contest.CaseResult{
			ID:       cid,
			Message:  r.Message,
			Verdict:  resolveVerdict(r.Verdict),
			Memory:   r.Memory,
			Duration: r.Duration,
		}
		stores.Submissions.UpdatePending(w.ctx, id, cr)
		d = cr
		break
	case finalResult:
		fv, p := getFinalVerdict(r)
		d = fResult{
			CompilerOutput: r.CompilerOutput,
			Verdict:        fv,
			Points:         p,
			MaxPoints:      r.MaxPoints,
		}
		w.commitToDb(id, stores.Submissions.GetPendingResults(w.ctx, id), r, fv, p)
		defer w.DestroySubscribers(id)
		break
	default:
		d = data
		break
	}
	if d != nil && len(subscribers) > 0 {
		for i := range subscribers {
			select {
			case subscribers[i] <- d:
			}
		}
	}
}

func (w *worker) DestroySubscribers(id uint32) {
	stores.Submissions.DeletePending(w.ctx, id)
	a, ok := w.sm.Pop(id)
	if !ok {
		return
	}
	for i := range a {
		close(a[i])
		a[i] = nil
	}
}

func (w *worker) Subscribe(id uint32) chan interface{} {
	c := make(chan interface{}, 1)
	subscribers, _ := w.sm.Get(id)
	w.sm.Set(id, append(subscribers, c))
	return c
}

func (w *worker) Unsubscribe(id uint32, c chan interface{}) {
	subscribers, ok := w.sm.Get(id)
	if !ok {
		return
	}
	close(c)
	w.sm.Set(id, slices.DeleteFunc(subscribers, func(rc chan interface{}) bool {
		return rc == c
	}))
}
