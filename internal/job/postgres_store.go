package job

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

func (s *PostgresStore) Create(job *Job) error {
	job.ID = uuid.New().String()
	job.State = JobPending
	job.RetryCount = 0
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	if job.MaxRetries <= 0 {
		job.MaxRetries = 3
	}

	if job.TimeoutSeconds <= 0 {
		job.TimeoutSeconds = 30
	}

	query := `
	INSERT INTO jobs (
		id,
		type,
		payload,
		state,
		retry_count,
		max_retries,
		timeout_seconds,
		created_at,
		updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := s.db.Exec(
		context.Background(),
		query,
		job.ID,
		job.Type,
		job.Payload,
		job.State,
		job.RetryCount,
		job.MaxRetries,
		job.TimeoutSeconds,
		job.CreatedAt,
		job.UpdatedAt,
	)

	return err
}

func (s *PostgresStore) Get(id string) (Job, error) {
	query := `
	SELECT
		id,
		type,
		payload,
		state,
		retry_count,
		max_retries,
		timeout_seconds,
		created_at,
		updated_at
	FROM jobs
	WHERE id = $1
	`

	var job Job

	err := s.db.QueryRow(context.Background(), query, id).Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.State,
		&job.RetryCount,
		&job.MaxRetries,
		&job.TimeoutSeconds,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		return Job{}, ErrJobNotFound
	}

	return job, nil
}

func (s *PostgresStore) List() ([]Job, error) {
	query := `
	SELECT
		id,
		type,
		payload,
		state,
		retry_count,
		max_retries,
		timeout_seconds,
		created_at,
		updated_at
	FROM jobs
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var j Job

		err := rows.Scan(
			&j.ID,
			&j.Type,
			&j.Payload,
			&j.State,
			&j.RetryCount,
			&j.MaxRetries,
			&j.TimeoutSeconds,
			&j.CreatedAt,
			&j.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (s *PostgresStore) Transition(id string, to JobState) error {
	query := `
	UPDATE jobs
	SET state = $1,
	    updated_at = $2
	WHERE id = $3
	`

	_, err := s.db.Exec(
		context.Background(),
		query,
		to,
		time.Now(),
		id,
	)

	return err
}

func (s *PostgresStore) IncrementRetry(id string) error {
	query := `
	UPDATE jobs
	SET retry_count = retry_count + 1,
	    updated_at = $1
	WHERE id = $2
	`

	_, err := s.db.Exec(
		context.Background(),
		query,
		time.Now(),
		id,
	)

	return err
}

func (s *PostgresStore) Pending() ([]Job, error) {
	query := `
	SELECT
		id,
		type,
		payload,
		state,
		retry_count,
		max_retries,
		timeout_seconds,
		created_at,
		updated_at
	FROM jobs
	WHERE state = $1
	`

	rows, err := s.db.Query(context.Background(), query, JobPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var j Job

		err := rows.Scan(
			&j.ID,
			&j.Type,
			&j.Payload,
			&j.State,
			&j.RetryCount,
			&j.MaxRetries,
			&j.TimeoutSeconds,
			&j.CreatedAt,
			&j.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}
