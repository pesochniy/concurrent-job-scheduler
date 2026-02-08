package job

import (
	"encoding/json"
	"time"
)

type JobState uint8

const (
	JobPending JobState = iota
	JobRunning
	JobCompleted
	JobFailed
	JobCanceled
)

type Job struct {
	ID      string
	Type    string
	Payload json.RawMessage

	State      JobState
	RetryCount int
	MaxRetries int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (j *Job) CanTransition(to JobState) bool {
	switch j.State {
	case JobPending:
		return to == JobRunning || to == JobCanceled
	case JobRunning:
		return to == JobCompleted || to == JobFailed || to == JobCanceled
	case JobFailed:
		return to == JobPending // retry
	case JobCompleted, JobCanceled:
		return false
	default:
		return false
	}
}
