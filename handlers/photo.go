package handlers

import (
	"context"
	"log"
	"time"

	"github.com/swastik-gautam/go-job-queue/job"
)

func resizePhoto(ctx context.Context, j job.Job) error {
	log.Printf("job %s: resizing photo, payload=%s", j.ID, j.Payload)

	select {
	case <-time.After(2 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
