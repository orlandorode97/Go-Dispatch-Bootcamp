package interactor

import (
	"errors"
	"strconv"
	"testing"

	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m MockDefinitionsRepo) GetConcurrentDefinitions(idType string, taskSize, perWorker int) (*model.List, error) {
	args := m.Called(idType, taskSize, perWorker)
	return args.Get(0).(*model.List), args.Error(1)
}

type MockDefinitionPresenter struct {
	mock.Mock
}

func (m MockDefinitionPresenter) ResponseDefinitions(definitionsList *model.List) (*model.List, error) {
	args := m.Called(definitionsList)
	return args.Get(0).(*model.List), args.Error(1)
}

func TestGet(t *testing.T) {
	testcases := []struct {
		name           string
		term           string
		response       *model.List
		responseLength int
		error          error
		hasError       bool
	}{
		{
			name: "success - response definitions by term",
			term: "the",
			response: &model.List{
				Definitions: []model.Definition{
					{
						Definition: "Some definition",
						Permalink:  "www.urbanctionary.example",
						ThumbsUp:   90,
						Author:     "Orlando",
						Word:       "the",
						Defid:      1234566,
						WrittenOn:  "2006-04-30T20:18:42.000Z",
						Example:    "[The] crusaders was a jazz-funk band founded at the beginning of the 70's ",
						ThumbsDown: 29,
					},
				},
			},
			responseLength: 1,
			error:          nil,
			hasError:       false,
		},
		{
			name:           "failure - definition was not found",
			term:           "a word that does not get records",
			responseLength: 0,
			error:          model.ErrNotFound{Term: "a word that does not get records"},
			hasError:       true,
		},
		{
			name:           "failure - json bad request",
			term:           "laksdasjlk lakjsdjajs lakjsda",
			responseLength: 0,
			error:          errors.New("unexpected end of JSON input"),
			hasError:       true,
		},
		{
			name:           "failure - missing api key",
			term:           "yo",
			responseLength: 0,
			error:          model.ErrMissingApiKey{},
			hasError:       true,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockDefinitionsRepo{}
			mockPresenter := MockDefinitionPresenter{}
			mockRepo.On("GetDefinitionsByTerm", test.term).Return(test.response, test.error)

			interactor := NewDefinitionInteractor(mockRepo, mockPresenter)
			definitions, err := interactor.urbanDictionaryRepository.GetDefinitionsByTerm(test.term)
			if test.hasError {
				assert.EqualError(t, err, test.error.Error())
			}
			if !test.hasError {
				assert.Nil(t, err)
			}
			if definitions != nil {
				assert.Equal(t, test.responseLength, len(definitions.Definitions))
			}
		})
	}
}

func TestGetFromCSV(t *testing.T) {
	testcases := []struct {
		name           string
		term           string
		response       *model.List
		definitionId   string
		error          error
		assertIdEquals func(assert.TestingT, interface{}, interface{}, ...interface{}) bool
	}{
		{
			name:         "success - definition found",
			term:         "yeah",
			definitionId: "1234566",
			response: &model.List{
				Definitions: []model.Definition{
					{
						Definition: "yeah definition",
						Permalink:  "www.urbanctionary.example",
						ThumbsUp:   90,
						Author:     "Orlando",
						Word:       "the",
						Defid:      1234566,
						WrittenOn:  "2006-04-30T20:18:42.000Z",
						Example:    "Hell [yeah]",
						ThumbsDown: 29,
					},
				},
			},
			assertIdEquals: assert.Equal,
		},
		{
			name:         "failure - definition not found",
			term:         "buddy",
			definitionId: "111111",
			response: &model.List{
				Definitions: []model.Definition{
					{
						Definition: "buddy definition",
						Permalink:  "www.urbanctionary.example",
						ThumbsUp:   90,
						Author:     "Orlando",
						Word:       "the",
						Defid:      1234566,
						WrittenOn:  "2006-04-30T20:18:42.000Z",
						Example:    "[The] movie",
						ThumbsDown: 29,
					},
				},
			},
			assertIdEquals: assert.NotEqual,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockDefinitionsRepo{}
			mockPresenter := MockDefinitionPresenter{}
			mockRepo.On("GetDefinitionById", test.definitionId).Return(test.response, test.error)

			interactor := NewDefinitionInteractor(mockRepo, mockPresenter)
			definitions, _ := interactor.urbanDictionaryRepository.GetDefinitionById(test.definitionId)
			id, err := strconv.Atoi(test.definitionId)
			assert.Nil(t, err)
			test.assertIdEquals(t, id, definitions.Definitions[0].Defid)
		})
	}
}
