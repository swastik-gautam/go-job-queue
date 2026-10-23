package store

import (
	"errors"
	"sync"

	"github.com/swastik-gautam/go-job-queue/job"
)

type Store struct {
	mu   sync.RWMutex
	jobs map[string]job
}

func NewStore() *Store {
	return &Store{
		jobs: make(map[string]job),
	}
}

func (s *Store) Add(j job.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[j.ID] = j
}

func (s *Store) SetState(id string, state job.Status) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	//find the existing job
	j, ok := s.job[id]
	if !ok {
		return errors.New("job not found")
	}

	//update its state
	j.State = state
	s.jobs[id] = j

	return nil
}

func (s *Store) Get(id string) (job.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	j, ok := s.jobs[id]
	if !ok {
		return job.Job{}, errors.New("job not found")
	}

	return j, nil
}
