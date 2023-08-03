package contest

import (
	"blizzard/blizzard/db/models/user"
	"github.com/google/uuid"
	"time"
)

type (
	Submission struct {
		// save the source code in a folder and load it by id
		ID          uint32    `bun:",pk,autoincrement" json:"id"`
		AuthorID    uuid.UUID `bun:",type:uuid"`
		ProblemID   string
		Language    string
		SubmittedAt *time.Time `bun:",nullzero,type:timestamptz,notnull,default:'now()'::timestamptz"`
		Verdict     Verdict    `bun:",nullzero" json:"verdict"`
		Author      *user.User
	}

	Verdict = uint8
)

const (
	Queued Verdict = iota
	Accepted
	InternalError
	Rejected
	Cancelled
	RuntimeError
	TimeLimitExceeded
	MemoryLimitExceeded
	OutputLimitExceeded
	StackLimitExceeded
	InvalidReturn
	CompilerError
)
