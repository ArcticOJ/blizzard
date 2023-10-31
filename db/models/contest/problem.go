package contest

import (
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/google/uuid"
)

type (
	Problem struct {
		ID            string       `bun:",pk" json:"id"`
		Tags          []string     `bun:",array" json:"tags"`
		Source        string       `json:"source"`
		AuthorID      uuid.UUID    `bun:",type:uuid" json:"authorID,omitempty"`
		Author        *user.User   `bun:"rel:has-one,join:author_id=id" json:"-"`
		Title         string       `bun:",notnull" json:"title,omitempty"`
		ContentType   ContentType  `json:"contentType"`
		Content       interface{}  `bun:"type:jsonb" json:"content,omitempty"`
		Constraints   *Constraints `bun:"embed:" json:"constraints,omitempty"`
		TestCount     uint16       `bun:",notnull" json:"testCount,omitempty"`
		PointsPerTest float64      `bun:",default:1" json:"pointPerTest,omitempty"`
		Submissions   []Submission `bun:"rel:has-many,join:id=problem_id" json:"submissions,omitempty"`
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

	ContentType = string
)

const (
	Structured ContentType = "structured"
	PDF                    = "pdf"
	Markdown               = "md"
	URL                    = "url"
	Image                  = "image"
)
