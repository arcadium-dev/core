package services_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/build"
	"arcadium.dev/core/http/services"
	"arcadium.dev/core/require"
)

func TestHealthRegister(t *testing.T) {
	method := http.MethodGet

	router := mux.NewRouter()
	s := services.Health{
		Start: time.Now(),
		Info:  build.Info("name", "version", "branch", "commit", "date"),
	}
	s.Register(router)

	r := httptest.NewRequest(method, services.HealthRoute, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	result := w.Result()

	assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
	assert.Equal(t, result.StatusCode, http.StatusOK)

	body, err := io.ReadAll(result.Body)
	require.Nil(t, err)
	require.NotEqual(t, len(body), 0)
	defer result.Body.Close()

	var resp services.HealthResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	assert.Equal(t, resp.Data.Name, "name")
	assert.Equal(t, resp.Data.Version, "version")
	assert.Equal(t, resp.Data.Build.Branch, "branch")
	assert.Equal(t, resp.Data.Build.Commit, "commit")
	assert.Equal(t, resp.Data.Build.Date, "date")
	assert.Equal(t, resp.Data.Build.Go, runtime.Version())
}

func TestHealthName(t *testing.T) {
	assert.Equal(t, services.Health{}.Name(), "health")
}
