package main

import (
	"fmt"

	"github.com/swastik-gautam/go-job-queue/job"
	"github.com/swastik-gautam/go-job-queue/producer"
	"github.com/swastik-gautam/go-job-queue/worker"
)

func main() {
	var n int

	fmt.Print("Enter number of jobs: ")
	if _, err := fmt.Scan(&n); err != nil || n < 0 {
		fmt.Println("Enter a valid non-negative number")
		return
	}

	// Create the shared buffered channel.
	jobs := make(chan job.Job, 100)

	// Start workers concurrently.
	workerCount := 3
	go worker.Start(jobs, workerCount)

	// Generate and submit jobs.
	producer.GenerateJobs(jobs, n)

	// Signal that no more jobs will be submitted.
	close(jobs)

	// TODO: Wait for workers to finish.
}
