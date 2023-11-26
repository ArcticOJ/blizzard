package judge

import (
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/polar/v0/types"
)

type (
	announcement struct {
		Type string `json:"type"`
		ID   uint16 `json:"id,omitempty"`
	}
	// final result for responding to clients
	fResult struct {
		CompilerOutput string          `json:"compilerOutput"`
		Verdict        contest.Verdict `json:"verdict"`
		Points         float64         `json:"points"`
		MaxPoints      float64         `json:"maxPoints"`
	}
)

func resolveVerdict(v types.CaseVerdict) contest.Verdict {
	switch v {
	case types.CaseVerdictAccepted:
		return contest.VerdictAccepted
	case types.CaseVerdictWrongAnswer:
		return contest.VerdictWrongAnswer
	case types.CaseVerdictInternalError:
		return contest.VerdictInternalError
	case types.CaseVerdictTimeLimitExceeded:
		return contest.VerdictTimeLimitExceeded
	case types.CaseVerdictMemoryLimitExceeded:
		return contest.VerdictMemoryLimitExceeded
	case types.CaseVerdictOutputLimitExceeded:
		return contest.VerdictOutputLimitExceeded
	case types.CaseVerdictRuntimeError:
		return contest.VerdictRuntimeError
	}
	return contest.VerdictNone
}
