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

func (s JobState) String() string {
	switch s {
	case JobPending:
		return "pending"
	case JobRunning:
		return "running"
	case JobCompleted:
		return "completed"
	case JobFailed:
		return "failed"
	case JobCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}
func (s JobState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type Job struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`

	State      JobState `json:"state"`
	RetryCount int      `json:"retry_count"`
	MaxRetries int      `json:"max_retries"`

	TimeoutSeconds int `json:"timeout_seconds"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
