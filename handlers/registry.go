package handlers

import (
	"context"
	"errors"

	"github.com/swastik-gautam/go-job-queue/job"
)

var ErrUnknownJobType = errors.New("unknown job type")

type HandlerFunc func(ctx context.Context, j job.Job) error

type Registry struct {
	m map[string]HandlerFunc
}

func NewRegistry() *Registry {
	return &Registry{m: make(map[string]HandlerFunc)}
}

func (r *Registry) Register(jobType string, h HandlerFunc) {
	r.m[jobType] = h
}

func (r *Registry) Lookup(jobType string) (HandlerFunc, bool) {
	h, ok := r.m[jobType]
	return h, ok
}
