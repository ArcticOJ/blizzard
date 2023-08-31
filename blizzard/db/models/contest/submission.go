package contest

import (
	"blizzard/blizzard/db/models/user"
	"github.com/google/uuid"
	"time"
)

type (
	Submission struct {
		ID       uint32    `bun:",pk,autoincrement" json:"id" json:"-"`
		AuthorID uuid.UUID `bun:",type:uuid" json:"-"`
		// file extension of the source code, we're using extension instead of full path because source code file name pattern is static except for the extension
		Extension   string     `bun:",notnull"`
		ProblemID   string     `json:"-"`
		Language    string     `json:"language"`
		SubmittedAt *time.Time `bun:",nullzero,type:timestamptz,notnull,default:'now()'::timestamptz" json:"submittedAt"`
		Result      *Result    `json:"result"`
		Author      *user.User `json:"author,omitempty"`
	}

	Verdict    string
	CaseResult struct {
		Message  string  `json:"message,omitempty"`
		Verdict  Verdict `json:"verdict"`
		Memory   uint32  `json:"memory"`
		Duration float32 `json:"duration"`
	}
	FinalResult struct {
		CompilerOutput string  `json:"compilerOutput"`
		Verdict        Verdict `json:"verdict"`
		Point          uint32  `json:"point"`
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
