package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	cache Cache
	mux   *http.ServeMux
}

func NewServer(cache Cache) *Server {
	s := &Server{cache: cache, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("POST /targeting/evaluate", s.handleEvaluate)
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
	Cached  bool `json:"cached"`
}

func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	var req evaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.Rule.FlagID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "rule.flag_id and user_id are required"})
		return
	}

	cacheKey := fmt.Sprintf("targeting:%s:%s", req.Rule.FlagID, req.UserID)

	if cached, found, err := s.cache.Get(r.Context(), cacheKey); err == nil && found {
		writeJSON(w, http.StatusOK, evaluateResponse{Matched: cached, Cached: true})
		return
	}

	matched := Evaluate(req.Rule, req.UserID)
	_ = s.cache.Set(r.Context(), cacheKey, matched)

	writeJSON(w, http.StatusOK, evaluateResponse{Matched: matched, Cached: false})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
