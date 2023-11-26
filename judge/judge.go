package judge

import (
	"github.com/ArcticOJ/blizzard/v0/db/models/contest"
	"github.com/ArcticOJ/polar/v0/types"
)

func getFinalVerdict(f types.FinalResult) (v contest.Verdict, points float64) {
	v = contest.VerdictNone
	points = f.Points
	if f.Verdict == types.FinalVerdictShortCircuit || f.Verdict == types.FinalVerdictNormal {
		v = contest.VerdictAccepted
		if f.LastNonACVerdict != types.CaseVerdictAccepted {
			v = resolveVerdict(f.LastNonACVerdict)
		}
		if f.Points > 0 && v != contest.VerdictAccepted {
			v = contest.VerdictPartiallyAccepted
		} else if v == contest.VerdictAccepted {
			points = f.MaxPoints
		}
	} else {
		switch f.Verdict {
		case types.FinalVerdictCancelled:
			v = contest.VerdictCancelled
		case types.FinalVerdictRejected:
			v = contest.VerdictRejected
		case types.FinalVerdictInitializationError:
			v = contest.VerdictInternalError
		case types.FinalCompileError:
			v = contest.VerdictCompileError
		}
	}
	return
}
