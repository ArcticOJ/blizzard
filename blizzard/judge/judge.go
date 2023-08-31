package judge

import (
	"blizzard/blizzard/db/models/contest"
)

type (
	Status struct {
		Name        string  `json:"name"`
		IsAlive     bool    `json:"isAlive"`
		Latency     float64 `json:"latency"`
		BootedSince int64   `json:"bootedSince"`
	}

	Submission struct {
		ID          uint32
		SourcePath  string
		Language    string
		ProblemID   string
		TestCount   uint16
		Constraints *contest.Constraints
	}
)

func resolveFinalResult(previousCases []contest.CaseResult, f FinalResult) *contest.FinalResult {
	fres := &contest.FinalResult{
		CompilerOutput: f.CompilerOutput,
		Verdict:        contest.None,
	}
	if f.Verdict == ShortCircuit || f.Verdict == Normal {
		var v contest.Verdict = contest.Accepted
		for i := range previousCases {
			if _v := previousCases[len(previousCases)-i-1].Verdict; _v != contest.Accepted {
				v = _v
				break
			}
		}
		fres.Verdict = v
	} else {
		v := contest.None
		switch f.Verdict {
		case Cancelled:
			v = contest.Cancelled
		case Rejected:
			v = contest.Rejected
		case InitError:
			v = contest.InternalError
		case CompileError:
			v = contest.CompilerError
		}
		fres.Verdict = v
	}
	return fres
}
