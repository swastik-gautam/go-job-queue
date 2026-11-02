package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/swastik-gautam/go-job-queue/handlers"
	"github.com/swastik-gautam/go-job-queue/job"
	"github.com/swastik-gautam/go-job-queue/store"
)

type Worker struct {
	store    *store.Store
	registry *handlers.Registry
}

func New(s *store.Store, r *handlers.Registry) *Worker {
	return &Worker{
		store:    s,
		registry: r,
	}
}

func Start(
	jobs <-chan job.Job,
	workerCount int,
	wg *sync.WaitGroup,
	w *Worker,
) {
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for j := range jobs {
				log.Printf("Worker %d received job %s", workerID, j.ID)
				w.execute(workerID, j)
			}
		}(i)
	}
}

func (w *Worker) execute(workerID int, j job.Job) {

	if err := w.store.SetState(j.ID, job.StatusRunning); err != nil {
		log.Printf("Worker %d: failed to update job state: %v", workerID, err)
		return
	}

	handler, ok := w.registry.Lookup(j.Type)

	if !ok {
		err := fmt.Errorf("unknown job type: %s", j.Type)
		w.fail(workerID, j, err)
		return
	}

	timeout := j.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	log.Printf("Worker %d executing job %s (%s)",
		workerID, j.ID, j.Type)

	if err := handler(ctx, j); err != nil {
		w.fail(workerID, j, err)
		return
	}

	if err := w.store.SetState(j.ID, job.StatusDone); err != nil {
		log.Printf("Worker %d: failed to mark job done: %v",
			workerID, err)
		return
	}

	log.Printf("Worker %d completed job %s", workerID, j.ID)
}

func (w *Worker) fail(workerID int, j job.Job, err error) {
	if stateErr := w.store.SetState(j.ID, job.StatusFailed); stateErr != nil {
		log.Printf("Worker %d: failed to update failed state: %v",
			workerID, stateErr)
	}

	log.Printf("Worker %d: job %s failed: %v",
		workerID, j.ID, err)
}
