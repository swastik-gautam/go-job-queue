package producer

import (
	"errors"

	"github.com/swastik-gautam/go-job-queue/job"
)

var ErrQueueFull = errors.New("queue is full")

func Enqueue(jobs chan<- job.Job, j job.Job) error {
	select {
	case jobs <- j:
		return nil
	default:
		return ErrQueueFull
	}
}
