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
		PointsPerTest float64
		Constraints   contest.Constraints
	}
)

func getFinalVerdict(f finalResult) (v contest.Verdict, points float64) {
	v = contest.None
	points = f.Points
	if f.Verdict == ShortCircuit || f.Verdict == Normal {
		v = contest.Accepted
		if f.LastNonACVerdict != contest.None {
			v = f.LastNonACVerdict
		}
		if f.Points > 0 && v != contest.Accepted {
			v = contest.PartiallyAccepted
		} else if v == contest.Accepted {
			points = f.MaxPoints
		}
	} else {
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
	}
	return
}
