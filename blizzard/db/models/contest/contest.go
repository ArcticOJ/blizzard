package contest

import (
	"blizzard/blizzard/db/models/user"
)

type (
	Contest struct {
		ID         uint32      `bun:",pk,autoincrement" json:"id"`
		Title      string      `bun:",notnull" json:"title"`
		Tags       []string    `bun:",array,notnull" json:"tags"`
		Organizers []user.User `bun:"m2m:contest_to_organizers,join:Contest=User"`
		Problems   []Problem   `bun:"m2m:contest_to_problems,join:Contest=Problem" json:"contests"`
	}

	ContestToOrganizer struct {
		ContestID uint64     `bun:",pk"`
		Contest   *Contest   `bun:"rel:belongs-to,join:contest_id=id"`
		UserID    string     `bun:",pk"`
		User      *user.User `bun:"rel:belongs-to,join:user_id=id"`
	}

	ContestToProblem struct {
		ProblemID string   `bun:",pk"`
		Problem   Problem  `bun:"rel:belongs-to,join:problem_id=id"`
		ContestID uint64   `bun:",pk"`
		Contest   *Contest `bun:"rel:belongs-to,join:contest_id=id"`
	}
)
