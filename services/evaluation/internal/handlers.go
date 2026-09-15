package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	publisher Publisher
	mux       *http.ServeMux
}

func NewServer(publisher Publisher) *Server {
	s := &Server{publisher: publisher, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("POST /evaluate", s.handleEvaluate)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type evaluateRequest struct {
	Rule   Rule   `json:"rule"`
	UserID string `json:"user_id"`
}

type evaluateResponse struct {
	Matched bool `json:"matched"`
}

func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	var req evaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.Rule.FlagID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "rule.flag_id and user_id are required"})
		return
	}

	matched := Evaluate(req.Rule, req.UserID)

	event := Event{
		FlagID:    req.Rule.FlagID,
		UserID:    req.UserID,
		Matched:   matched,
		Timestamp: time.Now().UTC(),
	}
	if err := s.publisher.Publish(r.Context(), event); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not publish analytics event"})
		return
	}

	writeJSON(w, http.StatusOK, evaluateResponse{Matched: matched})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
