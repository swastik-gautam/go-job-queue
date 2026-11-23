package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/swastik-gautam/go-job-queue/job"
	"github.com/swastik-gautam/go-job-queue/store"
)

// CreateJobRequest represents the JSON received from the client.
type CreateJobRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Handler handles incoming API requests.
type Handler struct {
	store *store.PostgresStore
}

// New creates a new API handler.
func New(s *store.PostgresStore) *Handler {
	return &Handler{
		store: s,
	}
}

// CreateJob validates the request and saves a new job in PostgreSQL.
func (h *Handler) CreateJob(c *gin.Context) {
	var req CreateJobRequest

	// Parse the incoming JSON request.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// Validate the job payload.
	payload := bytes.TrimSpace(req.Payload)

	if len(payload) == 0 || !json.Valid(payload) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "payload must be valid JSON",
		})
		return
	}

	// Only JSON objects are accepted as job payloads.
	if payload[0] != '{' {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "payload must be a valid JSON object",
		})
		return
	}

	// Validate the job type.
	if strings.TrimSpace(req.Type) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "type is required",
		})
		return
	}

	// Create a new job with its default values.
	now := time.Now().UTC()

	j := job.Job{
		ID:          uuid.NewString(),
		Type:        strings.TrimSpace(req.Type),
		Payload:     payload,
		Status:      job.StatusPending,
		Priority:    0,
		RunAt:       now,
		MaxAttempts: 3,
		Timeout:     30 * time.Second,
		CreatedAt:   now,
	}

	// Save the job in PostgreSQL.
	err := h.store.AddJob(c.Request.Context(), j)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to persist job",
		})
		return
	}

	// Return the job ID after successful persistence.
	c.JSON(http.StatusAccepted, gin.H{
		"message": "job accepted",
		"id":      j.ID,
	})
}
