package judge

import (
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/pb"
	"blizzard/blizzard/utils"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var ResultWatcher = New()

type (
	Watcher struct {
		fm cmap.ConcurrentMap[uint32, *feedback]
		sm map[uint32][]chan interface{}
		// sub channel
		sc chan *subscription
		// unsub channel
		usc chan *subscription
		// pub channel
		pc chan *message
		// dispose channel
		dc chan uint32
	}

	feedback struct {
		caseData []*pb.CaseResult
		final    *finalResult
	}

	finalResult struct {
		CompilerOutput string `yaml:"compilerOutput,omitempty"`
		Verdict        contest.Verdict
	}

	subscription struct {
		id uint32
		c  chan interface{}
	}

	message struct {
		id   uint32
		data interface{}
	}
)

func shardFunc(key uint32) uint32 {
	return key
}

func New() *Watcher {
	return &Watcher{
		fm:  cmap.NewWithCustomShardingFunction[uint32, *feedback](shardFunc),
		sm:  make(map[uint32][]chan interface{}),
		sc:  make(chan *subscription, 1),
		pc:  make(chan *message, 1),
		usc: make(chan *subscription, 1),
	}
}

func processFinalVerdict(metadata *pb.Metadata, r []*pb.JudgeResult) contest.Verdict {
	return contest.Accepted
}

func (w *Watcher) Watch() {
	for {
		select {
		// on dispose
		case id := <-w.dc:
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
			m = utils.ArrayRemove(m, func(r chan interface{}) bool {
				if r == u.c {
					u.c = nil
					return true
				}
				return false
			})
			if m == nil || len(m) == 0 {
				delete(w.sm, u.id)
			} else {
				w.sm[u.id] = m
			}
		// on pub
		case msg := <-w.pc:
			a, ok := w.sm[msg.id]
			if !ok {
				continue
			}
			for i := range a {
				select {
				case a[i] <- msg.data:
				}
			}
		}
	}
}

func (w *Watcher) addToQueue(sub *contest.Submission) {

}

func (w *Watcher) emit(id uint32, data interface{}) {
	w.pc <- &message{id: id, data: data}
}

func (w *Watcher) Queued() []*contest.Submission {
	return nil
}

func (w *Watcher) Dispose(id uint32) {
	w.dc <- id
}

func (w *Watcher) Track(id uint32, c chan interface{}) {
	w.sc <- &subscription{
		id: id,
		c:  c,
	}
}

func (w *Watcher) Untrack(id uint32, c chan interface{}) {
	w.usc <- &subscription{
		id: id,
		c:  c,
	}
}
