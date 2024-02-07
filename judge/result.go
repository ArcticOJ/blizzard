package judge

import (
	"github.com/ArcticOJ/blizzard/v0/db/schema/contest"
	"github.com/ArcticOJ/polar/v0/types"
)

const (
	typeManifest responseType = "manifest"
	typeAck      responseType = "ack"
	typeCase     responseType = "case"
	typeFinal    responseType = "final"
)

type (
	responseType = string
	// response with type to distinguish response types
	response struct {
		Type responseType `json:"type"`
		Data interface{}  `json:"data,omitempty"`
	}
	manifest struct {
		SubmissionID uint32  `json:"submissionId"`
		TestCount    uint16  `json:"testCount"`
		MaxPoints    float64 `json:"maxPoints"`
		// for initial payload
		AdditionalData interface{} `json:"additionalData,omitempty"`
	}
	// final judgement for responding to clients
	finalJudgement struct {
		CompilerOutput string          `json:"compilerOutput"`
		Verdict        contest.Verdict `json:"verdict"`
		Points         float64         `json:"points"`
	}
)

func resolveVerdict(v types.CaseVerdict) contest.Verdict {
	switch v {
	case types.CaseVerdictAccepted:
		return contest.VerdictAccepted
	case types.CaseVerdictWrongAnswer:
		return contest.VerdictWrongAnswer
	case types.CaseVerdictInternalError:
		return contest.VerdictInternalError
	case types.CaseVerdictTimeLimitExceeded:
		return contest.VerdictTimeLimitExceeded
	case types.CaseVerdictMemoryLimitExceeded:
		return contest.VerdictMemoryLimitExceeded
	case types.CaseVerdictOutputLimitExceeded:
		return contest.VerdictOutputLimitExceeded
	case types.CaseVerdictRuntimeError:
		return contest.VerdictRuntimeError
	}
	return contest.VerdictNone
}
