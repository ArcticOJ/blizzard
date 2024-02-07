package judge

import (
	"container/list"
	"context"
	"errors"
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/contest"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"github.com/ArcticOJ/polar/v0"
	"github.com/ArcticOJ/polar/v0/types"
	"github.com/google/uuid"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/mitchellh/mapstructure"
	"github.com/uptrace/bun"
	"math"
	"sync"
)

var Observer *observer

type (
	observer struct {
		ctx context.Context
		// subscribers
		m  sync.RWMutex
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
	}).Order("submitted_at DESC").Where("verdict = ?", "").Scan(ctx) != nil {
		return nil
	}
	r = make([]types.Submission, len(s))
	for i, _s := range s {
		c := _s.Problem.Constraints
		r[i] = types.Submission{
			AuthorID:      _s.AuthorID.String(),
			ID:            _s.ID,
			SourcePath:    _s.FileName,
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
	Observer = &observer{
		ctx: ctx,
		sm: csmap.Create[uint32, *submissionSubscribers](csmap.WithCustomHasher[uint32, *submissionSubscribers](func(key uint32) uint64 {
			return uint64(key)
		}), csmap.WithShardCount[uint32, *submissionSubscribers](1)),
	}
	Observer.p = polar.NewPolar(ctx, Observer.handleResult)
	Observer.p.Populate(getPendingSubmissions(ctx))
	Observer.p.StartServer()
}

func (o *observer) commitToDb(id uint32, cases []contest.CaseResult, fr types.FinalResult, v contest.Verdict, p float64) {
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
	if _, e := db.Database.NewUpdate().Model(&s).WherePK().Column("results", "verdict", "points", "compiler_output").Returning("NULL").Exec(o.ctx); e != nil {
		logger.Blizzard.Error().Err(e).Uint32("id", id).Msg("could not commit results to database")
	}
}

func (o *observer) Enqueue(sub types.Submission, subscribe bool) (chan interface{}, *list.Element, error) {
	var (
		c       chan interface{}
		element *list.Element
	)
	o.sm.Store(sub.ID, &submissionSubscribers{
		l: list.New(),
	})
	if subscribe {
		c, element = o.Subscribe(sub, func() interface{} {
			return nil
		})
	}
	return c, element, o.p.Push(sub, false)
}

func (o *observer) Cancel(id uint32, userId uuid.UUID) error {
	if !o.p.Cancel(id, userId.String()) {
		return errors.New("could not cancel specified submission")
	}
	return nil
}

func (o *observer) handleResult(id uint32) func(t types.ResultType, data interface{}) bool {
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
			o.publish(id, _r.CaseID, _r)
		case types.ResultFinal:
			var _r types.FinalResult
			if mapstructure.Decode(data, &_r) != nil {
				return false
			}
			_r.LastNonACVerdict = lastNonAcVerdict
			o.publish(id, math.MaxUint16, _r)
			return true
		case types.ResultAck:
			o.publish(id, math.MaxUint16, uint16(0))
		}
		return false
	}
}

func (o *observer) publish(id uint32, cid uint16, data interface{}) {
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
		o.p.UpdateResult(id, cr)
		d = &response{
			Type: typeCase,
			Data: cr,
		}
	case types.FinalResult:
		fv, p := getFinalVerdict(r)
		d = &response{
			Type: typeFinal,
			Data: finalJudgement{
				CompilerOutput: r.CompilerOutput,
				Verdict:        fv,
				Points:         p,
			},
		}
		o.commitToDb(id, o.p.GetResults(id), r, fv, p)
		defer o.DestroySubscribers(id)
	case uint16:
		d = &response{
			Type: typeAck,
		}
	}
	if d != nil {
		if subscribers, ok := o.sm.Load(id); ok {
			subscribers.m.RLock()
			for v := subscribers.l.Front(); v != nil; v = v.Next() {
				select {
				case v.Value.(chan interface{}) <- d:
				default:
				}
			}
			subscribers.m.RUnlock()
		}
	}
}

func (o *observer) DestroySubscribers(id uint32) {
	subscribers, ok := o.sm.Load(id)
	if !ok {
		return
	}
	subscribers.m.Lock()
	defer subscribers.m.Unlock()
	o.sm.Delete(id)
	// iterate over linked list and then close + delete all subscribers
	for v := subscribers.l.Front(); v != nil; v = v.Next() {
		close(v.Value.(chan interface{}))
	}
	subscribers.l.Init()
	subscribers = nil
}

func (o *observer) Subscribe(sub types.Submission, getData func() interface{}) (chan interface{}, *list.Element) {
	subscribers, ok := o.sm.Load(sub.ID)
	if !ok {
		return nil, nil
	}
	subscribers.m.Lock()
	// acknowledgement + n test case results + final result = n + 2
	c := make(chan interface{}, sub.TestCount+2)
	c <- response{
		Type: typeManifest,
		Data: manifest{
			SubmissionID:   sub.ID,
			TestCount:      sub.TestCount,
			MaxPoints:      float64(sub.TestCount) * sub.PointsPerTest,
			AdditionalData: getData(),
		},
	}
	defer subscribers.m.Unlock()
	return c, subscribers.l.PushBack(c)
}

func (o *observer) Unsubscribe(id uint32, e *list.Element) {
	subscribers, ok := o.sm.Load(id)
	if !ok {
		return
	}
	subscribers.m.Lock()
	defer subscribers.m.Unlock()
	close(e.Value.(chan interface{}))
	subscribers.l.Remove(e)
}

func (o *observer) RuntimeSupported(rt string) bool {
	return o.p.RuntimeAvailable(rt)
}

func (o *observer) GetJudges() map[string]*polar.JudgeObj {
	return o.p.GetJudges()
}

func (o *observer) GetResults(id uint32) []contest.CaseResult {
	return o.p.GetResults(id)
}
