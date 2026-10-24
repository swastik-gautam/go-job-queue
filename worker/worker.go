package worker

import (
	"log"
	"sync"
	"time"

	"github.com/swastik-gautam/go-job-queue/job"
	"github.com/swastik-gautam/go-job-queue/store"
)

type Worker struct {
	store *store.Store
}

func New(s *store.Store) *Worker {
	return &Worker{
		store: s,
	}
}

func Start(jobs <-chan job.Job, workerCount int, wg *sync.WaitGroup, s *store.Store) {
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for j := range jobs {
				s.SetState(j.ID, job.StatusRunning)

				log.Printf("Worker %d started job %s", workerID, j.ID)

				time.Sleep(500 * time.Millisecond)

				s.SetState(j.ID, job.StatusDone)

				log.Printf("Worker %d finished job %s", workerID, j.ID)
			}
		}(i)
	}
}
