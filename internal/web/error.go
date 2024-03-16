package web

import (
	"fmt"
	"net/http"
)

type errorCode string

type ErrorInterface interface {
	Code() string
	Description() string
	HTTPStatusCode() int
	Error() string
	Cause() string
}

type ErrorFunc func(desc string) ErrorInterface

const (
	UnauthorizedRequest = "unauthorized"
	BadRequest          = "bad_request"
	InternalServerError = "internal_server_error"
	NotFound            = "not_found"
)

var (
	ErrUnauthorizedRequest = func(desc string) ErrorInterface {
		return newError(UnauthorizedRequest, desc, "", http.StatusUnauthorized)
	}
	ErrBadRequest = func(desc string) ErrorInterface {
		return newError(BadRequest, desc, "", http.StatusBadRequest)
	}
	ErrInternalServerError = func(desc string) ErrorInterface {
		return newError(InternalServerError, desc, "", http.StatusInternalServerError)
	}
	ErrUpstreamError = func(code string, desc string) ErrorInterface {
		return newError(errorCode(code), desc, "", http.StatusOK)
	}
)

func newError(errCode errorCode, desc string, cause string, httpCode int) ErrorInterface {
	return &customError{code: errCode, description: desc, httpStatusCode: httpCode, cause: cause}
}

type customError struct {
	code           errorCode
	description    string
	cause          string
	httpStatusCode int
	metadata       any
}

func (e *customError) Code() string {
	return string(e.code)
}

func (e *customError) Description() string {
	return e.description
}

func (e *customError) HTTPStatusCode() int {
	return e.httpStatusCode
}

func (e *customError) Error() string {
	return fmt.Sprintf("code: %s description: %s httpStatusCode: %d cause: %s",
		e.code,
		e.description,
		e.httpStatusCode,
		e.cause,
	)
}

func (e *customError) Cause() string {
	return e.cause
}
