package judge

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/core/errs"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/ArcticOJ/blizzard/v0/storage"
	"github.com/ArcticOJ/polar/v0"
	"github.com/ArcticOJ/polar/v0/types"
	"github.com/google/uuid"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/mitchellh/mapstructure"
	"github.com/uptrace/bun"
	"io"
	"math"
	"sync"
)

var Worker *worker

type (
	worker struct {
		ctx context.Context
		// subscribers
		sm *csmap.CsMap[uint32, *submissionSubscribers]
		p  *polar.Polar
	}
	submissionSubscribers struct {
		m sync.RWMutex
		l *list.List
	}
)

func getPendingSubmissions(ctx context.Context) (r []types.Submission) {
	var s []contest.Submission
	if db.Database.NewSelect().Model(&s).Relation("Problem", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.ExcludeColumn("tags", "source", "author_id", "title", "content_type", "content")
	}).Order("submitted_at DESC").Where("results IS ?", nil).Scan(ctx) != nil {
		return nil
	}
	r = make([]types.Submission, len(s))
	for i, _s := range s {
		c := _s.Problem.Constraints
		r[i] = types.Submission{
			AuthorID:      _s.AuthorID.String(),
			ID:            _s.ID,
			SourcePath:    fmt.Sprintf("%d.%s", _s.ID, _s.Extension),
			Runtime:       _s.Runtime,
			ProblemID:     _s.ProblemID,
			TestCount:     _s.Problem.TestCount,
			PointsPerTest: _s.Problem.PointsPerTest,
			Constraints: types.Constraints{
				IsInteractive: c.IsInteractive,
				TimeLimit:     c.TimeLimit,
				MemoryLimit:   c.MemoryLimit,
				OutputLimit:   c.OutputLimit,
				AllowPartial:  c.AllowPartial,
				ShortCircuit:  c.ShortCircuit,
			},
		}
	}
	return
}

func Init(ctx context.Context) {
	Worker = &worker{
		ctx: ctx,
		sm:  csmap.Create[uint32, *submissionSubscribers](),
	}
	Worker.p = polar.NewPolar(ctx, Worker.handleResult)
	Worker.p.Populate(getPendingSubmissions(ctx))
	Worker.p.StartServer()
}

func (w *worker) commitToDb(id uint32, cases []contest.CaseResult, fr types.FinalResult, v contest.Verdict, p float64) {
	s := contest.Submission{
		ID:             id,
		Results:        cases,
		Verdict:        v,
		Points:         p,
		CompilerOutput: fr.CompilerOutput,
		TotalMemory:    0,
		TimeTaken:      0,
	}
	for _, r := range cases {
		s.TotalMemory += uint64(r.Memory)
		s.TimeTaken += r.Duration
	}
	logger.Panic(db.Database.RunInTx(w.ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, e := tx.NewUpdate().Model(&s).WherePK().Column("results", "verdict", "points", "compiler_output").Returning("NULL").Exec(w.ctx)
		return e
	}), "tx")
}

func (w *worker) Enqueue(sub types.Submission, subscribe bool, path string, f io.Reader) (chan interface{}, *list.Element, error) {
	var (
		c       chan interface{}
		element *list.Element
	)
	if !w.p.RuntimeAvailable(sub.Runtime) {
		return nil, nil, errs.JudgeNotAvailable
	}
	if e := storage.Submission.Write(path, f); e != nil {
		return nil, nil, e
	}
	w.sm.Store(sub.ID, &submissionSubscribers{
		l: list.New(),
	})
	if subscribe {
		c, element = w.Subscribe(sub.ID)
	}
	return c, element, w.p.Push(sub, false)
}

func (w *worker) Cancel(id uint32, userId uuid.UUID) error {
	if !w.p.Cancel(id, userId.String()) {
		return errors.New("could not cancel specified submission")
	}
	return nil
}

func (w *worker) handleResult(id uint32) func(t types.ResultType, data interface{}) bool {
	lastNonAcVerdict := types.CaseVerdictAccepted
	return func(t types.ResultType, data interface{}) bool {
		switch t {
		case types.ResultCase:
			var _r types.CaseResult
			if mapstructure.Decode(data, &_r) != nil {
				return false
			}
			if _r.Verdict != types.CaseVerdictAccepted {
				lastNonAcVerdict = _r.Verdict
			}
			w.publish(id, _r.CaseID, _r)
		case types.ResultFinal:
			var _r types.FinalResult
			if mapstructure.Decode(data, &_r) != nil {
				return false
			}
			_r.LastNonACVerdict = lastNonAcVerdict
			w.publish(id, math.MaxUint16, _r)
			return true
		case types.ResultAnnouncement:
			w.publish(id, math.MaxUint16, data.(uint16))
		}
		return false
	}
}

func (w *worker) publish(id uint32, cid uint16, data interface{}) {
	var d *response = nil
	switch r := data.(type) {
	case types.CaseResult:
		cr := contest.CaseResult{
			ID:       cid,
			Message:  r.Message,
			Verdict:  resolveVerdict(r.Verdict),
			Memory:   r.Memory,
			Duration: r.Duration,
		}
		w.p.UpdateResult(id, cr)
		d = &response{
			Type: typeCase,
			Data: cr,
		}
	case types.FinalResult:
		fv, p := getFinalVerdict(r)
		d = &response{
			Type: typeFinal,
			Data: fResult{
				CompilerOutput: r.CompilerOutput,
				Verdict:        fv,
				Points:         p,
				MaxPoints:      r.MaxPoints,
			},
		}
		w.commitToDb(id, w.p.GetResult(id), r, fv, p)
		defer w.DestroySubscribers(id)
	case uint16:
		d = &response{
			Type: typeAnnouncement,
			Data: data,
		}
	}
	if d != nil {
		subscribers, ok := w.sm.Load(id)
		if ok {
			subscribers.m.RLock()
			for v := subscribers.l.Front(); v != nil; v = v.Next() {
				select {
				case v.Value.(chan interface{}) <- d:
				}
			}
			subscribers.m.RUnlock()
		}
	}
}

func (w *worker) DestroySubscribers(id uint32) {
	subscribers, ok := w.sm.Load(id)
	if !ok {
		return
	}
	subscribers.m.Lock()
	w.sm.Delete(id)
	subscribers.m.Unlock()
	// iterate over linked list and then close + delete all subscribers
	for v := subscribers.l.Front(); v != nil; v = v.Next() {
		close(v.Value.(chan interface{}))
		subscribers.l.Remove(v)
	}
}

func (w *worker) Subscribe(id uint32) (chan interface{}, *list.Element) {
	subscribers, ok := w.sm.Load(id)
	if !ok {
		return nil, nil
	}
	subscribers.m.Lock()
	defer subscribers.m.Unlock()
	c := make(chan interface{}, 1)
	return c, subscribers.l.PushBack(c)
}

func (w *worker) Unsubscribe(id uint32, e *list.Element) {
	subscribers, ok := w.sm.Load(id)
	if !ok {
		return
	}
	subscribers.m.Lock()
	defer subscribers.m.Unlock()
	close(e.Value.(chan interface{}))
	subscribers.l.Remove(e)
}

func (w *worker) RuntimeSupported(rt string) bool {
	return w.p.RuntimeAvailable(rt)
}

func (w *worker) GetJudges() map[string]*polar.JudgeObj {
	return w.p.GetJudges()
}
