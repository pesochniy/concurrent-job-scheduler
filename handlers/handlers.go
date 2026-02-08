package handlers

import (
	"encoding/json"
	"net/http"
)

type JobRequest struct {
	JobType string          `json:"job_type"`
	Payload json.RawMessage `json:"payload"`
}

// Register registers HTTP routes on the provided mux.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/hello", helloHandler)
	mux.HandleFunc("POST /jobs", createJobHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello " + name})
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	w.WriteHeader(http.StatusOK)

}
