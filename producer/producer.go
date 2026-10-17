package producer

import (
	"github.com/swastik-gautam/go-job-queue/job"
)

func Enqueue(jobs chan<- job.Job, j job.Job) {
	jobs <- j
}
