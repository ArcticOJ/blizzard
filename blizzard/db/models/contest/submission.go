package contest

import "blizzard/blizzard/db/models/user"

type Submission struct {
	// save the source code in a folder and load it by id
	ID       uint64 `bun:",pk,autoincrement" json:"id"`
	AuthorID string
	Author   *user.User
}
