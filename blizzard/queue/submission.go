package queue

import (
	"blizzard/blizzard/pb/igloo"
	"sync"
)

type SubmissionQueue struct {
	resultChan chan *igloo.Submission
	wg         sync.WaitGroup
}

func (q *SubmissionQueue) Enqueue(submission *igloo.Submission) {
	q.Enqueue(submission)
}
