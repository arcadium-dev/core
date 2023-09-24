package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/http/middleware"
	"arcadium.dev/core/log"
)

func TestLoggingRequests(t *testing.T) {
	ctx, b := log.SetupTestLogging(t)

	router := mux.NewRouter()
	router.Use(middleware.Logging{Logger: zerolog.Ctx(ctx)}.Requests)
	router.HandleFunc("/foo", handler).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	assert.Equal(t, rw.Code, http.StatusOK)

	n := b.Len()
	assert.Contains(t, b.Index(n-3), `{"severity":"debug","method":"GET","url":"/foo","message":"request start"}`)
	assert.Contains(t, b.Index(n-2), `{"severity":"debug","method":"GET","url":"/foo","message":"handler called"}`)
	assert.Contains(t, b.Index(n-1), `{"severity":"debug","method":"GET","url":"/foo","message":"request complete"}`)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
	logger.Debug().Msg("handler called")
	w.WriteHeader(http.StatusOK)
}
