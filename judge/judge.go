package judge

import (
	"blizzard/db/models/contest"
)

type (
	Submission struct {
		ID          uint32
		SourcePath  string
		Language    string
		ProblemID   string
		TestCount   uint16
		Constraints *contest.Constraints
	}
)

func resolveFinalResult(f FinalResult) *contest.FinalResult {
	fres := &contest.FinalResult{
		CompilerOutput: f.CompilerOutput,
		Verdict:        contest.None,
	}
	if f.Verdict == ShortCircuit || f.Verdict == Normal {
		var v contest.Verdict = contest.Accepted
		if f.LastNonACVerdict != contest.None {
			v = f.LastNonACVerdict
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
