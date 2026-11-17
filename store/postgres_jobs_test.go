package store

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/swastik-gautam/go-job-queue/job"
)

func TestAddJob(t *testing.T) {
	ctx := context.Background()

	store, err := NewPostgresStore(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	j := job.Job{
		ID:      uuid.NewString(),
		Type:    "email",
		Payload: json.RawMessage(`{"to":"test@example.com"}`),
	}

	err = store.AddJob(ctx, j)
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}
}

func TestClaimJobs(t *testing.T) {
	ctx := context.Background()

	store, err := NewPostgresStore(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	jobs, err := store.ClaimJobs(ctx, 2)
	if err != nil {
		t.Fatalf("ClaimJobs failed: %v", err)
	}

	for _, j := range jobs {
		if j.Status != job.StatusRunning {
			t.Errorf("expected running, got %s", j.Status)
		}

		if j.Attempts < 1 {
			t.Errorf("expected attempts >= 1, got %d", j.Attempts)
		}
	}

	t.Logf("Successfully claimed %d jobs", len(jobs))
}
