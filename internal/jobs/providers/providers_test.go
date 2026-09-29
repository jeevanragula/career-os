package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeevanragula/career-os/internal/jobs"
)

func TestLeverAdapter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "application/json"); _, _ = w.Write([]byte(`[{"id":"abc","text":"Principal Engineer","categories":{"location":"Remote","commitment":"Full-time"},"descriptionPlain":"Build distributed systems","hostedUrl":"https://jobs.lever.co/acme/abc"}]`)) }))
	defer srv.Close()
	got, err := (LeverAdapter{Client: srv.Client()}).Discover(context.Background(), jobs.SourceConfig{Name:"acme", BaseURL:srv.URL}, jobs.DiscoveryQuery{})
	if err != nil { t.Fatal(err) }
	if len(got.Observations) != 1 || got.Observations[0].SourceJobID != "abc" { t.Fatalf("unexpected result: %+v", got) }
}

func TestAshbyAdapter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "application/json"); _, _ = w.Write([]byte(`{"jobs":[{"id":"xyz","title":"Staff Engineer","location":"Hyderabad","jobUrl":"https://jobs.ashbyhq.com/acme/xyz","descriptionPlain":"Build platforms"}]}`)) }))
	defer srv.Close()
	got, err := (AshbyAdapter{Client: srv.Client()}).Discover(context.Background(), jobs.SourceConfig{Name:"acme", BaseURL:srv.URL}, jobs.DiscoveryQuery{})
	if err != nil { t.Fatal(err) }
	if len(got.Observations) != 1 || got.Observations[0].SourceJobID != "xyz" { t.Fatalf("unexpected result: %+v", got) }
}
