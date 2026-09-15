package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRecordAndCounts(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()

	events := []Event{
		{FlagID: "flag-1", UserID: "u1", Matched: true, Timestamp: time.Now()},
		{FlagID: "flag-1", UserID: "u2", Matched: false, Timestamp: time.Now()},
		{FlagID: "flag-1", UserID: "u3", Matched: true, Timestamp: time.Now()},
	}

	for _, e := range events {
		if err := store.Record(ctx, e); err != nil {
			t.Fatalf("failed to record event: %v", err)
		}
	}

	matched, total, err := store.Counts(ctx, "flag-1")
	if err != nil {
		t.Fatalf("failed to get counts: %v", err)
	}
	if matched != 2 || total != 3 {
		t.Fatalf("expected matched=2 total=3, got matched=%d total=%d", matched, total)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	_ = store.Record(ctx, Event{FlagID: "flag-2", UserID: "u1", Matched: true, Timestamp: time.Now()})

	server := NewServer(store)
	req := httptest.NewRequest(http.MethodGet, "/metrics/flag-2", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp metricsResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Matched != 1 || resp.Total != 1 {
		t.Fatalf("expected matched=1 total=1, got matched=%d total=%d", resp.Matched, resp.Total)
	}
}

func TestMetricsEndpointUnknownFlag(t *testing.T) {
	server := NewServer(NewInMemoryStore())
	req := httptest.NewRequest(http.MethodGet, "/metrics/unknown", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp metricsResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Matched != 0 || resp.Total != 0 {
		t.Fatal("expected zeroed counts for unknown flag")
	}
}
