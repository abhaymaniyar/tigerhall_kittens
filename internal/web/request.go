package web

import (
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
	pathParams map[string]string
	params     map[string]string
}

type ValidationErrorInterface interface {
	Error() string
}

type ValidationError struct {
	errorType string
	message   string
	Err       error
}

func (e *ValidationError) Error() string { return e.message }

func NewRequest(r *http.Request) Request {
	return Request{Request: r}
}

func (r *Request) SetPathParam(key, value string) {
	if r.pathParams == nil {
		r.pathParams = make(map[string]string)
	}
	r.pathParams[key] = value
}

func (r *Request) GetPathParam(key string) string {
	if value, ok := r.pathParams[key]; ok {
		return value
	}
	return ""
}

func (r *Request) QueryParams() map[string]string {
	if r.params != nil {
		return r.params
	}
	r.params = map[string]string{}
	for key, val := range r.URL.Query() {
		r.params[key] = strings.Join(val, " | ")
	}
	return r.params
}
