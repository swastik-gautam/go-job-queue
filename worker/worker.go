package worker

import (
	"fmt"
	"sync"

	"github.com/swastik-gautam/go-job-queue/job"
)

func Start(jobs <-chan job.Job, workerCount int) {
	var wg sync.WaitGroup

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for j := range jobs {
				fmt.Printf("Worker %d processing job %s\n", workerID, j.ID)
			}
		}(i)
	}

	wg.Wait()
}
