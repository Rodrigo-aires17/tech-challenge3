package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAndGetFlag(t *testing.T) {
	server := NewServer(NewInMemoryRepository())

	body, _ := json.Marshal(createFlagRequest{Name: "new-checkout", Enabled: true})
	req := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var created Flag
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/flags/"+created.ID, nil)
	getRec := httptest.NewRecorder()
	server.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}
}

func TestGetFlagNotFound(t *testing.T) {
	server := NewServer(NewInMemoryRepository())

	req := httptest.NewRequest(http.MethodGet, "/flags/does-not-exist", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestUpdateAndDeleteFlag(t *testing.T) {
	repo := NewInMemoryRepository()
	server := NewServer(repo)

	created, err := repo.Create(context.Background(), "beta-banner", "", false)
	if err != nil {
		t.Fatalf("failed to seed flag: %v", err)
	}

	updateBody, _ := json.Marshal(updateFlagRequest{Enabled: true, Description: "rollout"})
	updateReq := httptest.NewRequest(http.MethodPut, "/flags/"+created.ID, strings.NewReader(string(updateBody)))
	updateRec := httptest.NewRecorder()
	server.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", updateRec.Code, updateRec.Body.String())
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/flags/"+created.ID, nil)
	deleteRec := httptest.NewRecorder()
	server.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", deleteRec.Code)
	}
}
