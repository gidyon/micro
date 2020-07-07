package http

import (
	"net/http"
)

// SupportCORS updates the CORs header for preflight requests
func SupportCORS(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Options request
		if r.Method == http.MethodOptions {
			w.Header().Set("access-control-allow-origin", "*")
			w.Header().Set("access-control-allow-methods", "POST, GET, PUT, PATCH, DELETE")
			w.Header().Set("access-control-allow-headers", "Authorization, Content-Type, Mode")
			return
		}

		// For Preflight request
		w.Header().Set("access-control-allow-origin", r.Header.Get("origin"))
		w.Header().Set("access-control-allow-credentials", "true")
		w.Header().Set("access-control-allow-headers", "Authorization, Content-Type, Mode")
		f.ServeHTTP(w, r)
	})
}
