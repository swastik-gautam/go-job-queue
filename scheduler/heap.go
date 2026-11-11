package scheduler

import (
	"container/heap"

	"github.com/swastik-gautam/go-job-queue/job"
)

type JobHeap []job.Job

func (h JobHeap) Len() int {
	return len(h)
}

func (h JobHeap) Less(i, j int) bool {
	return h[i].RunAt.Before(h[j].RunAt)
}

func (h JobHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *JobHeap) Push(x any) {
	*h = append(*h, x.(job.Job))
}

func (h *JobHeap) Pop() any {
	old := *h
	n := len(old)

	j := old[n-1]
	*h = old[:n-1]

	return j
}

func NewJobHeap() *JobHeap {
	h := &JobHeap{}
	heap.Init(h)
	return h
}
