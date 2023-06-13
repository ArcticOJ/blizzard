package problems

import (
	"github.com/uptrace/bun"
)

type Submission struct {
	bun.BaseModel
	// save the source code in a folder and load it by id
	ID uint64 `bun:",pk,autoincrement" json:"id"`
}
