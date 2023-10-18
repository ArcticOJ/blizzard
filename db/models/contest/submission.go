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
		Extension   string     `bun:",notnull" json:"extension"`
		ProblemID   string     `json:"-"`
		Problem     *Problem   `bun:"rel:belongs-to,join:problem_id=id" json:"problem,omitempty"`
		Language    string     `json:"language"`
		SubmittedAt *time.Time `bun:",nullzero,notnull,default:'now()'" json:"submittedAt"`
		Result      *Result    `json:"result"`
		Author      *user.User `bun:"rel:belongs-to,join:author_id=id" json:"author,omitempty"`
	}

	Verdict    string
	CaseResult struct {
		ID       uint16  `json:"id"`
		Message  string  `json:"message,omitempty"`
		Verdict  Verdict `json:"verdict"`
		Memory   uint32  `json:"memory"`
		Duration float32 `json:"duration"`
	}
	FinalResult struct {
		CompilerOutput string  `json:"compilerOutput"`
		Verdict        Verdict `json:"verdict"`
		Points         float32 `json:"points"`
		MaxPoints      float32 `json:"maxPoints"`
	}
	Result struct {
		Cases []CaseResult `json:"cases,omitempty"`
		Final *FinalResult `json:"final,omitempty"`
	}
)

const (
	None                Verdict = ""
	Accepted                    = "AC"
	PartiallyAccepted           = "PA"
	WrongAnswer                 = "WA"
	InternalError               = "IE"
	Rejected                    = "RJ"
	Cancelled                   = "CL"
	RuntimeError                = "RTE"
	TimeLimitExceeded           = "TLE"
	MemoryLimitExceeded         = "MLE"
	OutputLimitExceeded         = "OLE"
	StackLimitExceeded          = "SLE"
	CompilerError               = "CE"
)
