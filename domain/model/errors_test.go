package model

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		code        int
		contentType string
	}{
		{
			name:        "not found error type",
			err:         ErrNotFound{},
			code:        http.StatusNotFound,
			contentType: "application/json; charset=utf-8",
		},
		{
			name:        "not found in CSV error type",
			err:         ErrNotFoundInCSV{},
			code:        http.StatusNotFound,
			contentType: "application/json; charset=utf-8",
		},
		{
			name:        "missing api key error type",
			err:         ErrMissingApiKey{},
			code:        http.StatusForbidden,
			contentType: "application/json; charset=utf-8",
		},
		{
			name:        "missing api key secret error type",
			err:         ErrMissingApiKeySecret{},
			code:        http.StatusForbidden,
			contentType: "application/json; charset=utf-8",
		},
		{
			name:        "invalid data error type",
			err:         ErrInvalidData{},
			code:        http.StatusBadRequest,
			contentType: "application/json; charset=utf-8",
		},
		{
			name:        "random error type",
			err:         errors.New("random unknown error"),
			code:        http.StatusInternalServerError,
			contentType: "application/json; charset=utf-8",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			EncodeError(w, test.err)
			assert.Equal(t, test.code, w.Code)
			assert.Equal(t, test.contentType, w.Header().Get("Content-Type"))
		})
	}
}
