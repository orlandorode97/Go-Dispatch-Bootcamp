package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockInteractor struct {
	mock.Mock
}

func (m MockInteractor) Get(term string) (*model.List, error) {
	args := m.Called(term)
	return args.Get(0).(*model.List), args.Error(1)
}
func (m MockInteractor) GetFromCSV(id string) (*model.List, error) {
	args := m.Called(id)
	return args.Get(0).(*model.List), args.Error(1)
}
func (m MockInteractor) GetConcurrent(idType string, taskSize, perWorker int) (*model.List, error) {
	args := m.Called(idType, taskSize, perWorker)
	return args.Get(0).(*model.List), args.Error(1)
}

func TestGetDefinitions(t *testing.T) {
	testcases := []struct {
		name               string
		mockResponse       *model.List
		error              error
		term               string
		bodyMessage        string
		httpStatusExpected int
	}{
		{
			name:               "success - valid request",
			mockResponse:       &model.List{},
			error:              nil,
			term:               "sample",
			bodyMessage:        "the body response should be not nil",
			httpStatusExpected: http.StatusOK,
		},
		{
			name:               "failure - bad request",
			mockResponse:       &model.List{},
			error:              model.ErrInvalidData{Field: "term"},
			term:               "",
			bodyMessage:        "the body response should be not nil",
			httpStatusExpected: http.StatusBadRequest,
		},
		{
			name:               "success - definition not found",
			mockResponse:       &model.List{},
			error:              model.ErrNotFound{Term: "a random term"},
			term:               "a random term",
			bodyMessage:        "the body response should be not nil",
			httpStatusExpected: http.StatusNotFound,
		},
		{
			name:               "success - missing api key",
			mockResponse:       &model.List{},
			error:              model.ErrMissingApiKey{},
			term:               "wfio",
			bodyMessage:        "the body response should be not nil",
			httpStatusExpected: http.StatusForbidden,
		},
	}
	mockInteractor := MockInteractor{}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			mockInteractor.On("Get", test.term).Return(test.mockResponse, test.error)
			controller := NewDefinitionController(mockInteractor)
			req, err := http.NewRequest(http.MethodGet, "/definitions/", nil)
			assert.Nil(t, err, "new request error should be nil")
			req = mux.SetURLVars(req, map[string]string{"term": test.term})
			rec := httptest.NewRecorder()
			controller.GetDefinitions(rec, req)
			assert.NotNil(t, rec.Body, test.bodyMessage)
			assert.Equal(t, rec.Result().StatusCode, test.httpStatusExpected)
		})
	}
}
func TestGetDefinitionsFromCSV(t *testing.T) {

	testcases := []struct {
		name               string
		mockResponse       *model.List
		error              error
		id                 string
		bodyMessage        string
		httpStatusExpected int
		idExpected         int
	}{
		{
			name: "success - valid request by definition id",
			mockResponse: &model.List{
				Definitions: []model.Definition{
					{
						Word:  "hello",
						Defid: 1,
					},
				},
			},
			id:                 "1",
			error:              nil,
			bodyMessage:        "the body response should be not nil",
			httpStatusExpected: http.StatusOK,
			idExpected:         1,
		},
		{
			name:               "failure - definition not found",
			mockResponse:       &model.List{},
			id:                 "1233",
			error:              model.ErrNotFoundInCSV{Id: "1233"},
			bodyMessage:        "the body response should be not nil",
			httpStatusExpected: http.StatusNotFound,
		},
	}
	mockInteractor := MockInteractor{}
	var list model.List
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			mockInteractor.On("GetFromCSV", test.id).Return(test.mockResponse, test.error)
			controller := NewDefinitionController(mockInteractor)
			req, err := http.NewRequest(http.MethodGet, "/definitions/csv/", nil)
			assert.Nil(t, err, "new request error should be nil")
			req = mux.SetURLVars(req, map[string]string{"id": test.id})
			rec := httptest.NewRecorder()
			controller.GetDefinitionsFromCSV(rec, req)
			res := rec.Result()
			assert.NotNil(t, res.Body, test.bodyMessage)
			assert.Equal(t, res.StatusCode, test.httpStatusExpected)
			if test.idExpected != 0 {
				body, _ := ioutil.ReadAll(res.Body)
				json.Unmarshal(body, &list)
				assert.Equal(t, list.Definitions[0].Defid, test.idExpected)
			}
		})
	}
}
