package server_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/errors"
	"arcadium.dev/core/http/server"
)

func TestResponse(t *testing.T) {
	ctx := context.Background()

	t.Run("nil error", func(t *testing.T) {
		server.Response(context.Background(), nil, nil)
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			err  error
			code int
			body string
		}{
			{
				err:  errors.ErrBadRequest,
				code: http.StatusBadRequest,
				body: `{"status":400,"detail":"bad request"}` + "\n",
			},
			{
				err:  errors.ErrForbidden,
				code: http.StatusForbidden,
				body: `{"status":403,"detail":"forbidden"}` + "\n",
			},
			{
				err:  errors.ErrNotFound,
				code: http.StatusNotFound,
				body: `{"status":404,"detail":"not found"}` + "\n",
			},
			{
				err:  errors.ErrConflict,
				code: http.StatusConflict,
				body: `{"status":409,"detail":"conflict"}` + "\n",
			},
			{
				err:  errors.ErrInternal,
				code: http.StatusInternalServerError,
				body: `{"status":500,"detail":"internal server error"}` + "\n",
			},
			{
				err:  errors.ErrNotImplemented,
				code: http.StatusNotImplemented,
				body: `{"status":501,"detail":"not implemented"}` + "\n",
			},
			{
				err:  errors.New("unknown error"),
				code: http.StatusInternalServerError,
				body: `{"status":500,"detail":"unknown error"}` + "\n",
			},
			{
				err:  fmt.Errorf("failed to do something: %w", errors.ErrBadRequest),
				code: http.StatusBadRequest,
				body: `{"status":400,"detail":"failed to do something: bad request"}` + "\n",
			},
			{
				err:  fmt.Errorf("failed to authenticate: %w", errors.ErrForbidden),
				code: http.StatusForbidden,
				body: `{"status":403,"detail":"failed to authenticate: forbidden"}` + "\n",
			},
		}

		for _, test := range tests {
			w := httptest.NewRecorder()
			server.Response(ctx, w, test.err)

			assert.Equal(t, w.Header().Get("Content-Type"), "application/json; charset=utf-8")
			assert.Equal(t, w.Header().Get("X-Content-Type-Options"), "nosniff")
			assert.Equal(t, w.Code, test.code)
			assert.Equal(t, w.Body.String(), test.body)
		}
	})
}

func TestResponseError(t *testing.T) {
	expected := `status=200, detail="foobar"`
	err := server.ResponseError{Status: http.StatusOK, Detail: "foobar"}
	assert.Error(t, err, expected)
}
