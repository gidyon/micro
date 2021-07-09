package http

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"net/http/httputil"
// )

// // Example log output:
// // 127.0.0.1 - - [28/Oct/2016:18:35:05 -0400] "GET / HTTP/1.1" 200 13 "" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36"
// // 127.0.0.1 - - [28/Oct/2016:18:35:05 -0400] "GET /favicon.ico HTTP/1.1" 404 10 "http://localhost:8080/" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36"
// // Logs incoming requests, including response status.
// func Logger(h http.Handler, body bool) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		httputil.DumpRequest(r, body)
// 		h.ServeHTTP(w, r)
// 		httptest.Re
// 		res, ok := w.(*http.Response)
// 		if ok {
// 			httputil.DumpResponse(res, body)
// 		}
// 	})
// }
