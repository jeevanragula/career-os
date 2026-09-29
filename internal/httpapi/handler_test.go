package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()
	NewHandler().ServeHTTP(res, req)
	if res.Code != http.StatusOK { t.Fatalf("expected 200, got %d", res.Code) }
}
