package submissions

type (
	TestResult struct {
		Memory   int64   `json:"memory"`
		Duration int64   `json:"duration"`
		Verdict  Verdict `json:"verdict"`
		Point    int64   `json:"point"`
	}
	Verdict      byte
	FinalVerdict Verdict
)

const (
	Accepted Verdict = iota
	WrongAnswer
	TimeLimitExceeded
	MemoryLimitExceeded
	RuntimeError
)

const (
	PartiallyAccepted FinalVerdict = iota + 5
)
