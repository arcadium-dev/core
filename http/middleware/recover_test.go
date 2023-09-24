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

func TestRecoverPanics(t *testing.T) {
	ctx, b := log.SetupTestLogging(t)

	router := mux.NewRouter()
	router.Use(middleware.Recover{Logger: zerolog.Ctx(ctx)}.Panics)
	router.HandleFunc("/panic", boom).Methods(http.MethodPut)

	req := httptest.NewRequest(http.MethodPut, "/panic", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	assert.Equal(t, rw.Code, http.StatusInternalServerError)

	n := b.Len()
	assert.Contains(t, b.Index(n-2), `{"severity":"error","message":"recovering from a panic"}`)
	assert.Contains(t, b.Index(n-1), `{"severity":"error","message":"stacktrace:`)
}

func boom(http.ResponseWriter, *http.Request) {
	panic("boom")
}
