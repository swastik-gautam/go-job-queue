package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/swastik-gautam/go-job-queue/job"
)

func (s *PostgresStore) AddJob(ctx context.Context, j job.Job) error {
	// Validate required job fields
	if j.ID == "" {
		return fmt.Errorf("job ID cannot be empty")
	}

	if j.Type == "" {
		return fmt.Errorf("job type cannot be empty")
	}

	if !json.Valid(j.Payload) {
		return fmt.Errorf("invalid job payload")
	}

	// Set default values if not provided
	if j.Status == "" {
		j.Status = job.StatusPending
	}

	if j.RunAt.IsZero() {
		j.RunAt = time.Now()
	}

	if j.CreatedAt.IsZero() {
		j.CreatedAt = time.Now()
	}

	if j.MaxAttempts <= 0 {
		j.MaxAttempts = 3
	}

	if j.Timeout <= 0 {
		j.Timeout = 30 * time.Second
	}

	// Insert the job into PostgreSQL
	query := `
        INSERT INTO jobs (
            id,
            type,
            payload,
            status,
            priority,
            run_at,
            attempts,
            max_attempts,
            timeout_seconds,
            last_error,
            created_at
        )
        VALUES ($1, $2, $3::jsonb, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	_, err := s.DB.Exec(
		ctx,
		query,
		j.ID,
		j.Type,
		string(j.Payload),
		string(j.Status),
		j.Priority,
		j.RunAt,
		j.Attempts,
		j.MaxAttempts,
		int64(j.Timeout/time.Second),
		j.LastError,
		j.CreatedAt,
	)

	// Return an error if the database insert fails
	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}

	return nil

	if err != nil {
		fmt.Errorf("insert jobs %w", err)
	}

	return nil
}

func (s *PostgresStore) ClaimJobs(
	ctx context.Context,
	limit int,
) ([]job.Job, error) {

	if limit <= 0 {
		return nil, fmt.Errorf("claim limit must be greater than zero")
	}

	// Atomically select and claim available jobs
	query := `
        WITH picked AS (
            SELECT id
            FROM jobs
            WHERE status IN ('pending', 'scheduled')
              AND run_at <= NOW()
            ORDER BY priority DESC, run_at ASC, created_at ASC
            FOR UPDATE SKIP LOCKED
            LIMIT $1
        )
        UPDATE jobs AS j
        SET
            status = 'running',
            attempts = j.attempts + 1,
            started_at = NOW()
        FROM picked
        WHERE j.id = picked.id
        RETURNING
            j.id::text,
            j.type,
            j.payload,
            j.status,
            j.priority,
            j.run_at,
            j.attempts,
            j.max_attempts,
            j.timeout_seconds,
            j.last_error,
            j.created_at,
            j.started_at,
            j.finished_at
    `

	rows, err := s.DB.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("claim jobs: %w", err)
	}
	defer rows.Close()

	var claimedJobs []job.Job

	for rows.Next() {
		var j job.Job

		var status string
		var timeoutSeconds int32
		var lastError pgtype.Text
		var startedAt pgtype.Timestamptz
		var finishedAt pgtype.Timestamptz

		err := rows.Scan(
			&j.ID,
			&j.Type,
			&j.Payload,
			&status,
			&j.Priority,
			&j.RunAt,
			&j.Attempts,
			&j.MaxAttempts,
			&timeoutSeconds,
			&lastError,
			&j.CreatedAt,
			&startedAt,
			&finishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan claimed job: %w", err)
		}

		// Convert database values to Go types
		j.Status = job.Status(status)
		j.Timeout = time.Duration(timeoutSeconds) * time.Second

		if lastError.Valid {
			j.LastError = lastError.String
		}

		if startedAt.Valid {
			j.StartedAt = &startedAt.Time
		}

		if finishedAt.Valid {
			j.FinishedAt = &finishedAt.Time
		}

		claimedJobs = append(claimedJobs, j)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate claimed jobs: %w", err)
	}

	return claimedJobs, nil
}
