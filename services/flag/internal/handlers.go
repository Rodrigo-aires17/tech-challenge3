package internal

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Server struct {
	repo Repository
	mux  *http.ServeMux
}

func NewServer(repo Repository) *Server {
	s := &Server{repo: repo, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("GET /flags", s.handleList)
	s.mux.HandleFunc("POST /flags", s.handleCreate)
	s.mux.HandleFunc("GET /flags/{id}", s.handleGet)
	s.mux.HandleFunc("PUT /flags/{id}", s.handleUpdate)
	s.mux.HandleFunc("DELETE /flags/{id}", s.handleDelete)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type createFlagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req createFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	flag, err := s.repo.Create(r.Context(), req.Name, req.Description, req.Enabled)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create flag"})
		return
	}
	writeJSON(w, http.StatusCreated, flag)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	flags, err := s.repo.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not list flags"})
		return
	}
	writeJSON(w, http.StatusOK, flags)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	flag, err := s.repo.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "flag not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not get flag"})
		return
	}
	writeJSON(w, http.StatusOK, flag)
}

type updateFlagRequest struct {
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	flag, err := s.repo.Update(r.Context(), id, req.Enabled, req.Description)
	if errors.Is(err, ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "flag not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not update flag"})
		return
	}
	writeJSON(w, http.StatusOK, flag)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.repo.Delete(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "flag not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not delete flag"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
