package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swastik-gautam/go-job-queue/api"
	"github.com/swastik-gautam/go-job-queue/handlers"
	"github.com/swastik-gautam/go-job-queue/job"
	"github.com/swastik-gautam/go-job-queue/scheduler"
	"github.com/swastik-gautam/go-job-queue/store"
	"github.com/swastik-gautam/go-job-queue/worker"
)

func main() {
	// Shared channel used by all workers.
	jobs := make(chan job.Job, 100)

	jobStore := store.NewStore()

	var wg sync.WaitGroup

	// Register job handlers.
	registry := handlers.NewRegistry()
	handlers.RegisterAll(registry)

	// Create worker.
	w := worker.New(jobStore, registry)

	// Start three workers.
	worker.Start(jobs, 3, &wg, w)

	// Create and start delay scheduler.
	delayScheduler := scheduler.NewScheduler()
	delayScheduler.Start()

	// Forward expired jobs to the worker pool.
	go func() {
		for j := range delayScheduler.Ready() {
			jobs <- j
		}
	}()

	// Setup HTTP API.
	h := api.New(jobs, jobStore)

	r := gin.Default()
	r.POST("/jobs", h.CreateJob)
	r.GET("/jobs/:id", h.GetJob)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start HTTP server.
	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	stop()

	log.Println("shutting down...")

	// Shutdown HTTP server.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}

	// Close worker input and wait for workers.
	close(jobs)
	wg.Wait()

	log.Println("bye")
}
