package scheduler

import (
	"container/heap"
	"time"

	"github.com/swastik-gautam/go-job-queue/job"
)

type Scheduler struct {
	incoming chan job.Job
	ready    chan job.Job
	jobs     *JobHeap
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		incoming: make(chan job.Job),
		ready:    make(chan job.Job),
		jobs:     NewJobHeap(),
	}
}

func (s *Scheduler) Start() {
	go s.run()
}

func (s *Scheduler) Incoming() chan<- job.Job {
	return s.incoming
}

func (s *Scheduler) Ready() <-chan job.Job {
	return s.ready
}

func (s *Scheduler) run() {
	timer := time.NewTimer(time.Hour)
	timer.Stop()

	var timerC <-chan time.Time
	var pending *job.Job

	for {
		// Check whether we need to schedule a timer.
		if pending == nil && s.jobs.Len() > 0 {
			// Get the earliest job without removing it.
			next := (*s.jobs)[0]

			// Calculate how long until the job should run.
			delay := time.Until(next.RunAt)

			if delay < 0 {
				delay = 0
			}

			// Stop the previous timer and clean up if necessary.
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}

			// Start the timer for the next job.
			timer.Reset(delay)
			timerC = timer.C
		} else {
			// Disable the timer case when no timer is needed.
			timerC = nil
		}

		// Prepare the ready channel.
		readyC := s.ready
		var readyJob job.Job

		if pending != nil {
			readyJob = *pending
		} else {
			readyC = nil
		}

		// Wait for an event.
		select {
		case j := <-s.incoming:
			// Add the new job to the min heap.
			heap.Push(s.jobs, j)

		case <-timerC:
			// The earliest job's delay has expired.
			j := heap.Pop(s.jobs).(job.Job)
			pending = &j

		case readyC <- readyJob:
			// Successfully sent the job.
			pending = nil
		}
	}
}
