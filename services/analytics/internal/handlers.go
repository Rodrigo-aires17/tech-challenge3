package internal

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	store Store
	mux   *http.ServeMux
}

func NewServer(store Store) *Server {
	s := &Server{store: store, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /metrics/{flagID}", s.handleMetrics)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type metricsResponse struct {
	FlagID  string `json:"flag_id"`
	Matched int    `json:"matched"`
	Total   int    `json:"total"`
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	flagID := r.PathValue("flagID")
	matched, total, err := s.store.Counts(r.Context(), flagID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read metrics"})
		return
	}
	writeJSON(w, http.StatusOK, metricsResponse{FlagID: flagID, Matched: matched, Total: total})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
