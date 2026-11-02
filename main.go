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
	"github.com/swastik-gautam/go-job-queue/store"
	"github.com/swastik-gautam/go-job-queue/worker"
)

func main() {
	jobs := make(chan job.Job, 100)
	jobStore := store.NewStore()

	var wg sync.WaitGroup

	registry := handlers.NewRegistry()
	handlers.RegisterAll(registry)

	w := worker.New(jobStore, registry)

	worker.Start(jobs, 3, &wg, w)

	h := api.New(jobs, jobStore)
	r := gin.Default()
	r.POST("/jobs", h.CreateJob)
	r.GET("/jobs/:id", h.GetJob)

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM,
	)
	defer stop()
	<-ctx.Done()
	stop()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}

	close(jobs)
	wg.Wait()
	log.Println("bye")
}
