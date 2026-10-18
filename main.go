package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swastik-gautam/go-job-queue/handlers"
	"github.com/swastik-gautam/go-job-queue/job"
	"github.com/swastik-gautam/go-job-queue/worker"
)

func main() {
	jobs := make(chan job.Job, 100)

	worker.Start(jobs, 3)

	h := handlers.New(jobs)

	r := gin.Default()
	r.POST("/jobs", h.CreateJob)
	r.Run(":8080")
}
