package contest

import (
	"blizzard/blizzard/db/models/user"
	"github.com/google/uuid"
)

type (
	Problem struct {
		ID              string     `bun:",pk" json:"id"`
		Tags            []string   `bun:",array,notnull" json:"tags"`
		Source          string     `json:"source"`
		AuthorID        uuid.UUID  `bun:",type:uuid" json:"authorID,omitempty"`
		Author          *user.User `bun:"rel:has-one,join:author_id=id" json:"-"`
		*ProblemContent `bun:"embed:"`
		Constraints     *Constraints `bun:"embed:" json:"constraints,omitempty"`
		TestCount       uint16       `bun:",notnull" json:"testCount,omitempty"`
		PointPerTest    uint16       `bun:",default:1" json:"pointPerTest,omitempty"`
	}

	ProblemContent struct {
		Title           string            `bun:",notnull" json:"title,omitempty"`
		Statement       string            `bun:",notnull" json:"statement,omitempty"`
		Input           string            `bun:",notnull" json:"input,omitempty"`
		Output          string            `bun:",notnull" json:"output,omitempty"`
		Scoring         []string          `bun:",array" json:"scoring,omitempty"`
		SampleTestCases []SampleTestCases `json:"sampleTestCases,omitempty"`
	}

	SampleTestCases struct {
		Input  string `json:"input"`
		Output string `json:"output"`
		Note   string `json:"note"`
	}

	Constraints struct {
		IsInteractive    bool     `bun:",default:false" json:"isInteractive"`
		TimeLimit        float32  `bun:",default:1" json:"timeLimit"`
		MemoryLimit      uint32   `bun:",default:128" json:"memoryLimit"`
		OutputLimit      uint32   `bun:",default:64" json:"outputLimit"`
		AllowPartial     bool     `bun:",default:false" json:"allowPartial"`
		AllowedLanguages []string `bun:",array" json:"allowedLanguages"`
		ShortCircuit     bool     `bun:",default:false" json:"shortCircuit"`
	}
)
