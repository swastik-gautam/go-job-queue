package worker

import (
	"log"
	"sync"
	"time"

	"github.com/swastik-gautam/go-job-queue/job"
)

func Start(jobs <-chan job.Job, workerCount int, wg *sync.WaitGroup) {

	for i := 1; i <= workerCount; i++ {

		wg.Add(1)

		go func(workerID int) {

			for j := range jobs {

				log.Printf("Worker %d started job %s", workerID, j.ID)

				time.Sleep(500 * time.Millisecond)

				log.Printf("Worker %d Finished job %s OF (type: %s)", workerID, j.ID, j.Type)

				defer wg.Done()
			}
		}(i)
	}
}
