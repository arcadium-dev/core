// Copyright 2023 arcadium.dev <info@arcadium.dev>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package errors provides a set of tools for creating errors.
package errors // import "arcadium.dev/core/errors"

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	// Is is an alias for the errors.Is function. This allows this package to
	// used in place of the stdlib error package.
	Is = errors.Is

	// As is an alias for the errors.As function. This allows this package to
	// used in place of the stdlib error package.
	As = errors.As

	// New is an alias for the errors.New function. This allows this package to
	// used in place of the stdlib error package.
	New = errors.New
)

type (
	// HTTPError represents HTTP 4xx and 5xx classes of errors.
	HTTPError struct {
		status int
		msg    string
	}
)

// Status returns the http status code of the error.
func (e HTTPError) Status() int { return e.status }

// Error returns the error message.
func (e HTTPError) Error() string { return fmt.Sprintf("%s (Status: %d)", e.msg, e.status) }

var (
	// ErrBadRequest indicates an error due to client error.
	//
	// RFC 9110, 15.5.1
	ErrBadRequest = HTTPError{
		status: http.StatusBadRequest,
		msg:    strings.ToLower(http.StatusText(http.StatusBadRequest)),
	}

	// ErrForbidden indicates an authentication error due insufficient
	// credentials. The client should not automatically repeat the request with
	// the same credentials. The client MAY repeat the request with new or
	// different credentials.
	//
	// RFC 9110, 15.5.4
	ErrForbidden = HTTPError{
		status: http.StatusForbidden,
		msg:    strings.ToLower(http.StatusText(http.StatusForbidden)),
	}

	// ErrNotFound indicates that the requested resource was not found.
	//
	// RFC 9110, 15.5.5
	ErrNotFound = HTTPError{
		status: http.StatusNotFound,
		msg:    strings.ToLower(http.StatusText(http.StatusNotFound)),
	}

	// ErrConflict indicates that the current request could not be completed due
	// to the current status of the requested resource. This can occur on a POST
	// request where the resource already exists.
	//
	// RFC 9110, 15.5.10
	ErrConflict = HTTPError{
		status: http.StatusConflict,
		msg:    strings.ToLower(http.StatusText(http.StatusConflict)),
	}

	// ErrInternal indicates that the server encountered an unexpected condition
	// that prevented it from fulfilling the request.
	//
	// RFC 9110, 15.6.1
	ErrInternal = HTTPError{
		status: http.StatusInternalServerError,
		msg:    strings.ToLower(http.StatusText(http.StatusInternalServerError)),
	}

	// ErrNotImplemented  indicates that the server does not support the
	// functionality required to fulfill the request.
	//
	// // RFC 9110, 15.6.2
	ErrNotImplemented = HTTPError{
		status: http.StatusNotImplemented,
		msg:    strings.ToLower(http.StatusText(http.StatusNotImplemented)),
	}
)
