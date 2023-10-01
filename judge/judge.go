package judge

import (
	"blizzard/db/models/contest"
	"sync"
)

var l sync.RWMutex

var Status = make(map[string]*Judge)

func LockStatus() {
	l.Lock()
}

func UnlockStatus() {
	l.Unlock()
}

func GetStatus() map[string]*Judge {
	l.RLock()
	defer l.RUnlock()
	return Status
}

type (
	Judge struct {
		Alive    bool   `json:"alive"`
		Version  string `json:"version"`
		*Info    `json:"info"`
		Runtimes []Runtime `json:"runtimes"`
	}

	Info struct {
		Memory      uint32 `json:"memory"`
		OS          string `json:"os"`
		Parallelism uint8  `json:"parallelism"`
		BootedSince uint64 `json:"bootedSince"`
	}

	Runtime struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Compiler  string `json:"compiler"`
		Arguments string `json:"arguments"`
		Version   string `json:"version"`
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
