// Package http adds midlleware to http requests
package http

import (
	"context"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"time"
)

type id int

const key = id(12)

// AddRequestID adds request id to handler
func AddRequestID(h http.Handler) http.Handler {
	rand.Seed(time.Now().UnixNano())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, key, rand.Int()))
		h.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request id from context
func GetRequestID(ctx context.Context) (int, error) {
	v := ctx.Value(key)
	vint, ok := v.(id)
	if !ok {
		return 0, errors.New("failed to get value from context")
	}
	return int(vint), nil
}
