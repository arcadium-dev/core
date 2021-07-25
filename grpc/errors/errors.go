// Copyright 2021 arcadium.dev <info@arcadium.dev>
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

package errors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// Unwrap is an alias for errors.Unwrap
	Unwrap = errors.Unwrap

	// Is is an alias for errors.Is
	Is = errors.Is

	// As is an alias for errors.As
	As = errors.As
)

// New creates an error from the code and the message.
func New(code codes.Code, msg string) error {
	return &gerror{
		error: errors.New(msg),
		code:  code,
	}
}

// WithCode annotates the given error with the code.
func WithCode(err error, code codes.Code) error {
	return &gerror{
		error: err,
		code:  code,
	}
}

// Errorf creates an error from the code and a format specifier.
func Errorf(code codes.Code, format string, a ...interface{}) error {
	return New(code, fmt.Sprintf(format, a...))
}

// Wrap returns an error wrapped with the supplied message.
// If the given error is nil, Wrap returns nil.
func Wrap(e error, msg string) error {
	if e == nil {
		return nil
	}

	code := codes.Unknown
	ce, ok := e.(interface{ Code() codes.Code })
	if ok {
		code = ce.Code()
	}

	return &gerror{
		error: fmt.Errorf("%w: %s", e, msg),
		code:  code,
	}
}

// Wrap returns an error wrapped with the format specifier.
// If the given error is nil, Wrap returns nil.
func Wrapf(e error, format string, a ...interface{}) error {
	return Wrap(e, fmt.Sprintf(format, a...))
}

type (
	// gerror wraps a standard error with a grpc code.
	gerror struct {
		error
		code codes.Code
	}
)

// Unwrap returns an error if it was wrapped, and nil otherwise.
func (e *gerror) Unwrap() error {
	return errors.Unwrap(e.error)
}

// Core returns the grpc code associated with this error.
func (e *gerror) Code() codes.Code {
	return e.code
}

// GRPCStatus returns a grpc.Status, converting the game error
// to a status useful to be returned by a grpc service.
func (e *gerror) GRPCStatus() *status.Status {
	return status.New(e.code, e.Error())
}

var (
	// Base errors

	// Unknown indicates an unknown error.
	Unknown = New(codes.Unknown, codes.Unknown.String())

	// InvalidArgument indicates missing or malformed field in a request.
	InvalidArgument = New(codes.InvalidArgument, codes.InvalidArgument.String())

	// NotFound indicates that the requested resource does not exist.
	NotFound = New(codes.NotFound, codes.NotFound.String())

	// AlreadyExists indicates that the request to be created already exits.
	AlreadyExists = New(codes.AlreadyExists, codes.AlreadyExists.String())

	// Unimplemented indicates the requested service has not been implemented.
	Unimplemented = New(codes.Unimplemented, codes.Unimplemented.String())

	// Internal indicates an internal errors.
	Internal = New(codes.Internal, codes.Internal.String())
)
