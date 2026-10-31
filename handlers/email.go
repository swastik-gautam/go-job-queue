package handlers

import (
	"context"
	"log"

	"github.com/swastik-gautam/go-job-queue/job"
)

func sendEmail(ctx context.Context, j job.Job) error {
	log.Printf("job %s: sending email, payload=%s", j.ID, j.Payload)
	return nil
}
