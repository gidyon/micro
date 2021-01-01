package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

// CookieJWTOptions contains parameters for CookieToJWTMiddleware
type CookieJWTOptions struct {
	SecureCookie *securecookie.SecureCookie
	AuthHeader   string
	AuthScheme   string
	CookieName   string
}

// CookieToJWTMiddleware is a middleware for extracting cookie values to jwt
func CookieToJWTMiddleware(opt *CookieJWTOptions) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check that no authorization header is present
			if r.Header.Get(opt.AuthHeader) == "" {
				// Cookie name must be present
				cookie, err := r.Cookie(opt.CookieName)
				if err == nil {
					var token string
					// Decode cookie to token
					err = opt.SecureCookie.Decode(opt.CookieName, cookie.Value, &token)
					if err == nil {
						// Set authorization header
						r.Header.Set(opt.AuthHeader, fmt.Sprintf("Bearer %s", token))
					}
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
