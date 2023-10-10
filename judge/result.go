package judge

import "github.com/ArcticOJ/blizzard/v0/db/models/contest"

type (
	CaseVerdict  int8
	FinalVerdict int8
	CaseResult   struct {
		Duration float32
		Memory   uint32
		Message  string
		Verdict  CaseVerdict
	}
	FinalResult struct {
		CompilerOutput   string
		Verdict          FinalVerdict
		LastNonACVerdict contest.Verdict
	}
)

const (
	Accepted CaseVerdict = iota
	WrongAnswer
	InternalError
	TimeLimitExceeded
	MemoryLimitExceeded
	OutputLimitExceeded
	RuntimeError
)

const (
	Normal FinalVerdict = iota
	ShortCircuit
	Rejected
	Cancelled
	CompileError
	InitError
)

func resolveVerdict(v CaseVerdict) contest.Verdict {
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
