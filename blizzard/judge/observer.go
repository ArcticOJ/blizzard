package judge

import (
	"blizzard/blizzard/cache"
	"blizzard/blizzard/config"
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/logger"
	"blizzard/blizzard/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/uptrace/bun"
	"github.com/vmihailenco/msgpack"
	"math"
	"strconv"
)

var ResponseObserver = New()

type (
	Observer struct {
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
		rmq    amqp.Queue
		mqChan *amqp.Channel
		mqConn *amqp.Connection
	}

	Subscription struct {
		id uint32
		c  chan<- interface{}
	}

	result struct {
		id     uint32
		caseId uint16
		ttl    uint16
		data   interface{}
	}
)

func New() *Observer {
	// TODO: auto reconnect queue once closed
	conf := config.Config.RabbitMQ
	conn, e := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", conf.Username, conf.Password, conf.Host, conf.Port))
	logger.Panic(e, "failed to establish a connection to rabbitmq")
	ch, e := conn.Channel()
	logger.Panic(e, "failed to open a channel for queue")
	logger.Panic(ch.ExchangeDeclare("submissions", "direct", true, false, false, false, amqp.Table{
		"x-queue-type":       "quorum",
		"x-consumer-timeout": 3600000,
	}), "failed to declare exchange for submissions")
	rq, e := ch.QueueDeclare("results", false, false, false, false, amqp.Table{
		"x-single-active-consumer": true,
	})
	logger.Panic(e, "failed to open queue for results")
	logger.Panic(ch.Confirm(false), "could not put queue channel to confirm mode")
	return &Observer{
		sm:     make(map[uint32][]chan<- interface{}),
		sc:     make(chan Subscription, 1),
		pc:     make(chan result, 1),
		usc:    make(chan Subscription, 1),
		dc:     make(chan uint32, 1),
		rmq:    rq,
		mqChan: ch,
		mqConn: conn,
	}
}

func commitToDb(id uint32, res *contest.FinalResult) {
	db.Database.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		if _, e := tx.NewUpdate().Model((*contest.Submission)(nil)).Where("id = ?", id).Set("result = ?", res).Returning("NULL").Exec(context.Background()); e != nil {
			return e
		}
		return nil
	})
}

func (w *Observer) Enqueue(sub *Submission) error {
	b, e := msgpack.Marshal(sub)
	if e != nil {
		return e
	}
	e = cache.Result.Create(sub.ID, sub.TestCount)
	if e != nil {
		return e
	}
	conf, e := w.mqChan.PublishWithDeferredConfirmWithContext(
		context.Background(),
		"submissions",
		sub.Language,
		true,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/msgpack",
			Body:         b,
		})
	if e != nil {
		return e
	}
	return cache.Result.SetTag(sub.ID, conf.DeliveryTag)
}

func (w *Observer) Cancel(id uint32) error {
	if !cache.Result.IsPending(id) {
		return errors.New("no submission with matching ID")
	}
	if tag, e := cache.Result.GetTag(id); e == nil {
		fmt.Println(tag)
		return w.mqChan.Reject(tag, false)
	}
	return errors.New("could not cancel specified submission")
}

func (w *Observer) CreateResultConsumer() {
	c, e := w.mqChan.Consume(w.rmq.Name, "blizzard", true, false, false, false, nil)
	logger.Panic(e, "failed to start consuming judge results")
	for jr := range c {
		sid, e := strconv.ParseUint(fmt.Sprint(jr.Headers["submission-id"]), 10, 32)
		if e != nil {
			continue
		}
		if fmt.Sprint(jr.Headers["type"]) == "final" {
			var r FinalResult
			if msgpack.Unmarshal(jr.Body, &r) != nil {
				continue
			}
			w.publish(uint32(sid), math.MaxUint16, r, 0)
		} else if cid, e := strconv.ParseUint(fmt.Sprint(jr.Headers["case-id"]), 10, 16); e == nil {
			var r CaseResult
			var ttl uint16 = 0
			if _ttl, e := strconv.ParseUint(fmt.Sprint(jr.Headers["ttl"]), 10, 16); e == nil && _ttl > 0 {
				ttl = uint16(_ttl)
			}
			if msgpack.Unmarshal(jr.Body, &r) != nil {
				continue
			}
			w.publish(uint32(sid), uint16(cid), r, ttl)
		}
	}
}

func (w *Observer) Work() {
	go w.CreateResultConsumer()
	for {
		select {
		// on destroy
		case id := <-w.dc:
			cache.Result.Clean(id)
			cache.Result.DeleteTag(id)
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
			switch res := msg.data.(type) {
			case CaseResult:
				cr := contest.CaseResult{
					Message:  res.Message,
					Verdict:  resolveVerdict(res.Verdict),
					Memory:   res.Memory,
					Duration: res.Duration,
				}
				d = cr
				// TODO: use higher priority for case results to ensure race condition when delivering
				// TODO: handle cache failure
				cache.Result.Store(msg.id, msg.caseId, cr, int(msg.ttl))
			case FinalResult:
				r, _ := cache.Result.Get(msg.id)
				var cases []contest.CaseResult
				_ = json.Unmarshal([]byte(r), &cases)
				fr := resolveFinalResult(cases, res)
				d = fr
				commitToDb(msg.id, fr)
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

func (w *Observer) publish(id uint32, caseId uint16, data interface{}, ttl uint16) {
	w.pc <- result{id: id, caseId: caseId, data: data, ttl: ttl}
}

func (w *Observer) DestroyObserver(id uint32) {
	w.dc <- id
}

func (w *Observer) Observe(id uint32, c chan<- interface{}) (s Subscription) {
	s = Subscription{
		id: id,
		c:  c,
	}
	w.sc <- s
	return
}

func (w *Observer) StopObserve(s Subscription) {
	w.usc <- s
}
