package judge

import (
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
)

type (
	Submission struct {
		ID            uint32
		SourcePath    string
		Language      string
		ProblemID     string
		TestCount     uint16
		PointsPerTest float32
		Constraints   contest.Constraints
	}
)

func resolveFinalResult(f FinalResult) *contest.FinalResult {
	fr := &contest.FinalResult{
		CompilerOutput: f.CompilerOutput,
		Verdict:        contest.None,
		Points:         f.Points,
		MaxPoints:      f.MaxPoints,
	}
	if f.Verdict == ShortCircuit || f.Verdict == Normal {
		var v contest.Verdict = contest.Accepted
		if f.LastNonACVerdict != contest.None {
			v = f.LastNonACVerdict
		}
		if f.Points > 0 && v != contest.Accepted {
			v = contest.PartiallyAccepted
		} else if v == contest.Accepted {
			fr.Points = f.MaxPoints
		}
		fr.Verdict = v
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
		fr.Verdict = v
	}
	return fr
}
