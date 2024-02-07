package contest

import (
	"github.com/ArcticOJ/blizzard/v0/db/schema/user"
	"github.com/google/uuid"
	"time"
)

type (
	Submission struct {
		ID             uint32       `bun:",pk,autoincrement" json:"id" json:"-"`
		AuthorID       uuid.UUID    `bun:",type:uuid" json:"-"`
		FileName       string       `bun:",notnull" json:"-"`
		ProblemID      string       `json:"problemId"`
		Problem        *Problem     `bun:",rel:belongs-to,join:problem_id=id" json:"problem,omitempty"`
		Runtime        string       `json:"runtime"`
		SubmittedAt    time.Time    `bun:",nullzero,notnull,default:current_timestamp" json:"submittedAt"`
		Results        []CaseResult `json:"results"`
		TimeTaken      float32      `json:"timeTaken"`
		TotalMemory    uint64       `json:"totalMemory"`
		Verdict        Verdict      `json:"verdict"`
		Points         float64      `json:"points"`
		CompilerOutput string       `json:"compilerOutput"`
		Author         *user.User   `bun:",rel:belongs-to,join:author_id=id" json:"author,omitempty"`
	}

	Verdict    = string
	CaseResult struct {
		ID       uint16  `json:"id"`
		Message  string  `json:"message,omitempty"`
		Verdict  Verdict `json:"verdict"`
		Memory   uint32  `json:"memory"`
		Duration float32 `json:"duration"`
	}
)

const (
	VerdictNone                Verdict = ""
	VerdictAccepted            Verdict = "AC"
	VerdictPartiallyAccepted   Verdict = "PA"
	VerdictWrongAnswer         Verdict = "WA"
	VerdictInternalError       Verdict = "IE"
	VerdictRejected            Verdict = "RJ"
	VerdictCancelled           Verdict = "CL"
	VerdictRuntimeError        Verdict = "RTE"
	VerdictTimeLimitExceeded   Verdict = "TLE"
	VerdictMemoryLimitExceeded Verdict = "MLE"
	VerdictOutputLimitExceeded Verdict = "OLE"
	VerdictStackLimitExceeded  Verdict = "SLE"
	VerdictCompileError        Verdict = "CE"
)
