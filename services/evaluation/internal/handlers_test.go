package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEvaluateEndpointPublishesEvent(t *testing.T) {
	publisher := &InMemoryPublisher{}
	server := NewServer(publisher)

	body, _ := json.Marshal(evaluateRequest{
		Rule:   Rule{FlagID: "flag-1", RolloutPercentage: 100},
		UserID: "user-1",
	})
	req := httptest.NewRequest(http.MethodPost, "/evaluate", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp evaluateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Matched {
		t.Fatal("expected matched=true for 100% rollout")
	}

	if len(publisher.Events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(publisher.Events))
	}
	if publisher.Events[0].FlagID != "flag-1" || publisher.Events[0].UserID != "user-1" {
		t.Fatal("published event does not match request")
	}
}

func TestEvaluateEndpointValidatesRequest(t *testing.T) {
	server := NewServer(&InMemoryPublisher{})

	req := httptest.NewRequest(http.MethodPost, "/evaluate", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
