package services_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	services "arcadium.dev/core/http/services"
)

func TestMetricsRegister(t *testing.T) {
	method := http.MethodGet

	router := mux.NewRouter()
	s := services.Metrics{}
	s.Register(router)

	r := httptest.NewRequest(method, services.MetricsRoute, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body")
	}
	defer resp.Body.Close()

	if len(body) == 0 {
		t.Error("Expected a response body")
	}
}

func TestMetricsName(t *testing.T) {
	var s services.Metrics
	if s.Name() != "metrics" {
		t.Errorf("Unexpected service name: %s", s.Name())
	}
}
