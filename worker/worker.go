package worker

import (
	"log"
	"time"

	"github.com/swastik-gautam/go-job-queue/job"
)

func Start(jobs <-chan job.Job, workerCount int) {

	for i := 1; i <= workerCount; i++ {
		go func(workerID int) {
			for j := range jobs {
				log.Printf("Worker %d started job %s", workerID, j.ID)
				time.Sleep(500 * time.Millisecond)
				log.Printf("Worker %d Finished job %s OF (type: %s)", workerID, j.ID, j.Type)
			}
		}(i)
	}
}
