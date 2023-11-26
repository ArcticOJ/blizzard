package contest

import (
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/google/uuid"
	"time"
)

type (
	Submission struct {
		ID       uint32    `bun:",pk,autoincrement" json:"id" json:"-"`
		AuthorID uuid.UUID `bun:",type:uuid" json:"-"`
		// file extension of the source code, we're using extension instead of full path because source code file name pattern is static except for the extension
		Extension      string       `bun:",notnull" json:"extension"`
		ProblemID      string       `json:"problemId"`
		Problem        *Problem     `bun:"rel:belongs-to,join:problem_id=id" json:"problem,omitempty"`
		Runtime        string       `json:"runtime"`
		SubmittedAt    time.Time    `bun:",nullzero,notnull,default:'now()'" json:"submittedAt"`
		Results        []CaseResult `json:"results"`
		TimeTaken      float32      `json:"timeTaken"`
		TotalMemory    uint64       `json:"totalMemory"`
		Verdict        Verdict      `json:"verdict"`
		Points         float64      `json:"points"`
		CompilerOutput string       `json:"compilerOutput"`
		Author         *user.User   `bun:"rel:belongs-to,join:author_id=id" json:"author,omitempty"`
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
	VerdictAccepted                    = "AC"
	VerdictPartiallyAccepted           = "PA"
	VerdictWrongAnswer                 = "WA"
	VerdictInternalError               = "IE"
	VerdictRejected                    = "RJ"
	VerdictCancelled                   = "CL"
	VerdictRuntimeError                = "RTE"
	VerdictTimeLimitExceeded           = "TLE"
	VerdictMemoryLimitExceeded         = "MLE"
	VerdictOutputLimitExceeded         = "OLE"
	VerdictStackLimitExceeded          = "SLE"
	VerdictCompileError                = "CE"
)
