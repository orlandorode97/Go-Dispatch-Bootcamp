package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrInvalidData struct {
	Field string
}

func (e ErrInvalidData) Error() string {
	return fmt.Sprintf("the field `%s` is invalid", e.Field)
}

type ErrNotFound struct {
	ImageUrl string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("there are not words in the image %s", e.ImageUrl)
}

type ErrNotFoundInCSV struct {
	Id string
}

func (e ErrNotFoundInCSV) Error() string {
	return fmt.Sprintf("there is word with id %s", e.Id)
}

type ErrMissingApiKey struct{}

func (e ErrMissingApiKey) Error() string {
	return "invalid api key"
}

type ErrMissingApiKeySecret struct{}

func (e ErrMissingApiKeySecret) Error() string {
	return "invalid api key secret"
}

// ErrInvalidDataType returned when we can not correctly incode struct
type ErrInvalidDataType struct {
	InvalidExpected string
}

func (e ErrInvalidDataType) Error() string {
	m := "invalid data type"
	if e.InvalidExpected != "" {
		return fmt.Sprintf("%s '%s'", m, e.InvalidExpected)
	}
	return m
}

// EncodeError encodes the error into a json format and writing the corresponding http status
func EncodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err.(type) {
	case ErrNotFound, ErrNotFoundInCSV:
		w.WriteHeader(http.StatusNotFound)
	case ErrMissingApiKey, ErrMissingApiKeySecret:
		w.WriteHeader(http.StatusForbidden)
	case ErrInvalidData, ErrInvalidDataType:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
