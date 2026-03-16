package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pesochniy/concurrent-job-scheduler/internal/domain"
	jobpkg "github.com/pesochniy/concurrent-job-scheduler/internal/job"
	"github.com/pesochniy/concurrent-job-scheduler/internal/workers"
)

type Scheduler interface {
	Start(ctx context.Context) error
	Submit(id string)
	Stop() error
}

type SchedulerImpl struct {
	store   jobpkg.Store
	queue   chan string
	workers int
}

func NewScheduler(store jobpkg.Store, workers int) *SchedulerImpl {
	return &SchedulerImpl{
		store:   store,
		queue:   make(chan string, 100),
		workers: workers,
	}
}

func (s *SchedulerImpl) Submit(id string) {
	s.queue <- id
}

func (s *SchedulerImpl) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case id, ok := <-s.queue:
			if !ok {
				return
			}

			job, err := s.store.Get(id)
			if err != nil {
				continue
			}

			if err := s.store.Transition(id, jobpkg.JobRunning); err != nil {
				continue
			}

			timeout := job.TimeoutSeconds
			if timeout <= 0 {
				timeout = 30
			}

			ctxExec, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)

			err = s.executeJob(ctxExec, job)
			cancel()

			if err != nil {

				if err := s.store.IncrementRetry(id); err != nil {
					_ = s.store.Transition(id, jobpkg.JobFailed)
					continue
				}

				job.RetryCount++

				if job.MaxRetries > 0 && job.RetryCount >= job.MaxRetries {
					_ = s.store.Transition(id, jobpkg.JobFailed)
					continue
				}

				_ = s.store.Transition(id, jobpkg.JobPending)

				select {
				case s.queue <- id:
				default:
				}

				continue
			}

			_ = s.store.Transition(id, jobpkg.JobCompleted)
		}
	}
}

func (s *SchedulerImpl) Start(ctx context.Context) error {

	for i := 0; i < s.workers; i++ {
		go s.worker(ctx)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return nil

		case <-ticker.C:
			jobs, err := s.store.Pending()
			if err != nil {
				continue
			}

			for _, j := range jobs {
				s.Submit(j.ID)
			}
		}
	}
}

func (s *SchedulerImpl) Stop() error {
	close(s.queue)
	return nil
}

func (s *SchedulerImpl) executeJob(ctx context.Context, j jobpkg.Job) error {
	errCh := make(chan error, 1)

	go func() {
		switch j.Type {

		case "email":
			var payload domain.Email
			if err := json.Unmarshal(j.Payload, &payload); err != nil {
				errCh <- err
				return
			}
			errCh <- workers.SendEmail(payload.To, payload.Subject, payload.Body)

		case "fetch_url":
			var payload domain.FetchURL
			if err := json.Unmarshal(j.Payload, &payload); err != nil {
				errCh <- err
				return
			}
			_, err := workers.FetchURL(payload.URL, payload.Timeout)
			errCh <- err

		case "report":
			var payload domain.Report
			if err := json.Unmarshal(j.Payload, &payload); err != nil {
				errCh <- err
				return
			}
			_, err := workers.GenerateReport(payload.UserID, payload.From, payload.To)
			errCh <- err

		default:
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-errCh:
		return err
	}
}
