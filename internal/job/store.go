package job

import (
	"errors"
)

var (
	ErrJobNotFound        = errors.New("job not found")
	ErrInvalidTransition  = errors.New("invalid job state transition")
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
)

type Store interface {
	Create(job *Job) error
	Get(id string) (Job, error)
	Transition(id string, to JobState) error
	IncrementRetry(id string) error
	List() ([]Job, error)
	Pending() ([]Job, error)
}
