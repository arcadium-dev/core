package services_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/http/services"
)

func TestPProfRegister(t *testing.T) {
	method := http.MethodGet

	router := mux.NewRouter()
	s := services.PProf{}
	s.Register(router)

	r := httptest.NewRequest(method, services.PProfRoute, nil)
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, r)
	result := w.Result()

	assert.Equal(t, result.StatusCode, http.StatusOK)
}

func TestPProfName(t *testing.T) {
	assert.Equal(t, services.PProf{}.Name(), "pprof")
}
