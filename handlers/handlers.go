package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swastik-gautam/go-job-queue/job"
	"github.com/swastik-gautam/go-job-queue/producer"
)

type CreateJobRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Handler struct {
	jobs chan<- job.Job
}

func New(ch chan<- job.Job) *Handler {
	return &Handler{jobs: ch}
}

func (h *Handler) CreateJob(c *gin.Context) {
	var req CreateJobRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": " invalid request",
		})
		return
	}

	payload := bytes.TrimSpace(req.Payload)

	if len(payload) == 0 || !json.Valid(payload) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "payload must be valid",
		})
		return
	}

	if bytes.Equal(payload, []byte("null")) || payload[0] != '{' {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "payload must be valid json object",
		})
		return
	}

	if strings.TrimSpace(req.Type) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "type is required",
		})
		return
	}

	j := job.Job{
		ID:        uuid.NewString(),
		Type:      strings.TrimSpace(req.Type),
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
	}

	producer.Enqueue(h.jobs, j)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "request accepted",
		"id":      j.ID,
	})
}
