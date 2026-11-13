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
	defer timer.Stop()

	var timerC <-chan time.Time
	var pending *job.Job

	for {
		// Check the earliest job and configure timer
		if pending == nil && s.jobs.Len() > 0 {
			next := (*s.jobs)[0].RunAt
			delay := time.Until(next)

			if delay < 0 {
				delay = 0
			}

			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}

			timer.Reset(delay)
			timerC = timer.C
		} else {
			timerC = nil
		}

		// Prepare job delivery
		var readyC chan job.Job
		var readyJob job.Job

		if pending != nil {
			readyC = s.ready
			readyJob = *pending
		}

		// Wait for an event
		select {
		case j := <-s.incoming:
			heap.Push(s.jobs, j)

		case <-timerC:
			j := heap.Pop(s.jobs).(job.Job)
			pending = &j

		case readyC <- readyJob:
			pending = nil
		}
	}
}
