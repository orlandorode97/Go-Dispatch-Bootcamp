package middleware

import "net/http"

// HeadersMiddleware set the Content-Type for json responses
func HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header().Get("Content-Type")
		if header == "" {
			w.Header().Add("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}
