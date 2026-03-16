package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pesochniy/concurrent-job-scheduler/internal/job"
	"github.com/pesochniy/concurrent-job-scheduler/internal/scheduler"
)

type Handler struct {
	store     job.Store
	scheduler scheduler.Scheduler
}

func NewHandler(store job.Store, scheduler scheduler.Scheduler) *Handler {
	return &Handler{
		store:     store,
		scheduler: scheduler,
	}
}

type createJobRequest struct {
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
	MaxRetries     int             `json:"max_retries"`
	TimeoutSeconds int             `json:"timeout_seconds"`
}

// Register registers HTTP routes on the provided mux.
func Register(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/health", h.healthHandler)
	mux.HandleFunc("/api/hello", h.helloHandler)
	mux.HandleFunc("POST /jobs", h.createJobHandler)
	mux.HandleFunc("GET /jobs", h.listJobsHandler)
	mux.HandleFunc("GET /jobs/{id}", h.getJobHandler)
}

func (h *Handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello " + name})
}

func (h *Handler) createJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	var req createJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}
	switch req.Type {
	case "email", "fetch_url", "report":
		// valid
	default:
		http.Error(w, "unknown job type", http.StatusBadRequest)
		return
	}
	j := &job.Job{
		Type:           req.Type,
		Payload:        req.Payload,
		MaxRetries:     req.MaxRetries,
		TimeoutSeconds: req.TimeoutSeconds,
	}
	if err := h.store.Create(j); err != nil {
		http.Error(w, "failed to create job", http.StatusInternalServerError)
		return
	}
	if h.scheduler != nil {
		h.scheduler.Submit(j.ID)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"id": j.ID,
	})

}

func (h *Handler) getJobHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	j, err := h.store.Get(id)
	if err != nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(j)
}

func (h *Handler) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.store.List()
	if err != nil {
		http.Error(w, "failed to list jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jobs)
}
