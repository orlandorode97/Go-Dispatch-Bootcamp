package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHeadersMiddleware(t *testing.T) {
	t.Run("set up header Content-Type:application/json", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := mux.NewRouter()
		r.Use(HeadersMiddleware)
		r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {})
		req := httptest.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)
		jsonHeader := w.Header().Get("Content-Type")
		assert.Equal(t, "application/json", jsonHeader, "middleware should set up Content-Type:application/json")
	})
}
