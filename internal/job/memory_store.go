package job

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStore struct {
	mu   sync.Mutex
	jobs map[string]Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs: make(map[string]Job),
	}
}

func (s *MemoryStore) Create(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()

	job.ID = id
	job.State = JobPending
	job.RetryCount = 0
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	// Store a copy
	s.jobs[id] = *job

	return nil
}

func (s *MemoryStore) Get(id string) (Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return Job{}, ErrJobNotFound
	}

	return job, nil
}

func (s *MemoryStore) Transition(id string, to JobState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return ErrJobNotFound
	}

	if !job.CanTransition(to) {
		return ErrInvalidTransition
	}

	job.State = to
	job.UpdatedAt = time.Now()

	s.jobs[id] = job
	return nil
}

func (s *MemoryStore) IncrementRetry(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return ErrJobNotFound
	}

	if job.RetryCount >= job.MaxRetries {
		return ErrMaxRetriesExceeded
	}

	job.RetryCount++
	job.UpdatedAt = time.Now()

	s.jobs[id] = job
	return nil
}

func (s *MemoryStore) List() ([]Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobs := make([]Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		jobs = append(jobs, j)
	}

	return jobs, nil
}
