package job

import (
	"encoding/json"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusScheduled Status = "scheduled"
	StatusRunning   Status = "running"
	StatusDone      Status = "done"
	StatusFailed    Status = "failed"
	StatusDead      Status = "dead"
)

type Job struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`

	Status   Status `json:"status"`
	Priority int    `json:"priority"`

	RunAt time.Time `json:"run_at"`

	Attempts    int           `json:"attempts"`
	MaxAttempts int           `json:"max_attempts"`
	Timeout     time.Duration `json:"timeout"`
	LastError   string        `json:"last_error,omitempty"`

	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}
