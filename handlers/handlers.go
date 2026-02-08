package handlers

import (
	"encoding/json"
	"net/http"
)

// Register registers HTTP routes on the provided mux.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/hello", helloHandler)
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
