package utils

import (
	"net/http"

	"github.com/ggicci/httpin"
)

func FromRequestContext[T any](r *http.Request, key any) (T, bool) {
	var zero T

	if r == nil {
		return zero, false
	}

	v := r.Context().Value(key)
	if v == nil {
		return zero, false
	}

	// Case 1: stored as T
	if val, ok := v.(T); ok {
		return val, true
	}

	// Case 2: stored as *T
	if ptr, ok := v.(*T); ok && ptr != nil {
		return *ptr, true
	}

	return zero, false
}

func InputFromContext[T any](r *http.Request) (T, bool) {
	return FromRequestContext[T](r, httpin.Input)
}
