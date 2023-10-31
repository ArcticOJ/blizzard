package judge

import "github.com/ArcticOJ/blizzard/v0/db/models/contest"

type (
	caseVerdict  int8
	finalVerdict int8
	announcement struct {
		Type string `json:"type"`
		ID   uint16 `json:"id,omitempty"`
	}
	caseResult struct {
		Duration float32
		Memory   uint32
		Message  string
		Verdict  caseVerdict
	}
	// final result for deserializing response from judge
	finalResult struct {
		CompilerOutput   string
		Verdict          finalVerdict
		Points           float64
		MaxPoints        float64
		LastNonACVerdict contest.Verdict
	}
	// final result for responding to clients
	fResult struct {
		CompilerOutput string          `json:"compilerOutput"`
		Verdict        contest.Verdict `json:"verdict"`
		Points         float64         `json:"points"`
		MaxPoints      float64         `json:"maxPoints"`
	}
)

const (
	Accepted caseVerdict = iota
	WrongAnswer
	InternalError
	TimeLimitExceeded
	MemoryLimitExceeded
	OutputLimitExceeded
	RuntimeError
)

const (
	Normal finalVerdict = iota
	ShortCircuit
	Rejected
	Cancelled
	CompileError
	InitError
)

func resolveVerdict(v caseVerdict) contest.Verdict {
	switch v {
	case Accepted:
		return contest.Accepted
	case WrongAnswer:
		return contest.WrongAnswer
	case InternalError:
		return contest.InternalError
	case TimeLimitExceeded:
		return contest.TimeLimitExceeded
	case MemoryLimitExceeded:
		return contest.MemoryLimitExceeded
	case OutputLimitExceeded:
		return contest.OutputLimitExceeded
	case RuntimeError:
		return contest.RuntimeError
	}
	return contest.None
}
