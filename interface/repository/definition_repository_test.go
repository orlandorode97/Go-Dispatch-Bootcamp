package repository

import (
	"strconv"
	"testing"

	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/orlandorode97/go-disptach/infrastructure/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testApiKey = "cc2464e4a1mshb5ceeca91e5a6adp1fa80bjsn4b48e2408b87"
)

type MockDefinitionsRepo struct {
	mock.Mock
}

func (m MockDefinitionsRepo) GetDefinitionsByTerm(term string) (*model.List, error) {
	args := m.Called(term)
	return args.Get(0).(*model.List), args.Error(1)
}

func (m MockDefinitionsRepo) GetDefinitionById(id string) (*model.List, error) {
	args := m.Called(id)
	return args.Get(0).(*model.List), args.Error(1)
}

func (m MockDefinitionsRepo) GetConcurrentDefinitions(id string) (*model.List, error) {
	args := m.Called(id)
	return args.Get(0).(*model.List), args.Error(1)
}

func TestGetDefinitionsByTerm(t *testing.T) {
	testcases := []struct {
		name     string
		term     string
		response *model.List
		error    error
		hasError bool
	}{
		{
			name:     "success - get definitions by term",
			term:     "Yeeah",
			error:    nil,
			hasError: false,
		},
		{
			name:     "failure - definition not found",
			term:     "a",
			error:    model.ErrNotFound{Term: "a"},
			hasError: true,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			mock := MockDefinitionsRepo{}
			mock.On("GetDefinitionsByTerm", test.term).Return(test.response, test.error)
			urbanClient := api.NewUrbanDictionary(testApiKey)

			repo := NewUrbanDictionaryRepository(urbanClient)

			_, err := repo.urbanDictionaryClient.GetDefinitions(test.term)
			if test.hasError {
				assert.EqualError(t, err, test.error.Error())
			}
		})
	}
}

func TestGetDefinitionById(t *testing.T) {
	testcases := []struct {
		name         string
		term         string
		response     *model.List
		error        error
		definitionId string
		hasError     bool
	}{
		{
			name:         "success - get definition by id",
			term:         "Yeeah",
			definitionId: "10593002",
			error:        nil,
			hasError:     false,
		},
		{
			name:         "failure - definition not found",
			term:         "a",
			definitionId: "123456",
			error:        model.ErrNotFoundInCSV{Id: "123456"},
			hasError:     true,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			mock := MockDefinitionsRepo{}
			mock.On("GetDefinitionById", test.term).Return(test.response, test.error)
			urbanClient := api.NewUrbanDictionary(testApiKey)
			repo := NewUrbanDictionaryRepository(urbanClient)

			repo.urbanDictionaryClient.UpdateCSVPath("../../data/definitions_test.csv")

			definitions, err := repo.urbanDictionaryClient.GetDefinitionById(test.definitionId)
			if test.hasError {
				assert.EqualError(t, err, test.error.Error())
			}
			if !test.hasError {
				id := strconv.Itoa(definitions.Definitions[0].Defid)
				assert.Equal(t, test.definitionId, id)
			}
		})
	}
}
