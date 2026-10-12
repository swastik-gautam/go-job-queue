package producer

import (
	"fmt"
	"time"

	"github.com/swastik-gautam/go-job-queue/job"
)

func GenerateJobs(jobs chan<- job.Job, n int) {

	for i := 0; i < n; i++ {
		j := job.Job{
			ID:        fmt.Sprint("jobs-%d", i+1),
			Type:      "example",
			Payload:   []byte(fmt.Sprint("task-%d", i+1)),
			CreatedAt: time.Now(),
		}
		jobs <- j
	}
}
