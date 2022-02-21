package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/stretchr/testify/assert"
)

var (
	testApiKey       = "cc2464e4a1mshb5ceeca91e5a6adp1fa80bjsn4b48e2408b87"
	testApiKeySecret = "cc2464e4a1mshb5ceeca91e5a6adp1fa80bjsn4b48e2408b87"
)

type mockDoFunc func(req *http.Request) (*http.Response, error)

type MockClient struct {
	DoFunc mockDoFunc
}

func (m MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestUpdateCSVPath(t *testing.T) {
	t.Run("set valid CSV path", func(t *testing.T) {
		imaggaClient := New(testApiKey, testApiKeySecret)
		imaggaClient.UpdateCSVPath("../../data/words_valid.csv")
		assert.Equal(t, imaggaClient.CSVPath, "../../data/words_valid.csv")
	})
}

func TestGetDefinitions(t *testing.T) {
	successBody := `{
		"list": [
		  {
			"definition": "A term used by a [parent] meaning you. It is commonly used when they know you [will not] [do something] by yourself.",
			"permalink": "http://we.urbanup.com/4971898",
			"thumbs_up": 95,
			"author": "deadeye10000",
			"word": "we",
			"defid": 4971898,
			"written_on": "Mon, 17-May-2010 21:13",
			"example": "Father: hey son, lets go, we are going to go [mow the lawn].\r\nSon: sure, [sounds fun]!\r\nFather: you go ahead and start cutting, I will be out [in 5 minutes].\r\nSon: dripping from sweat, where were you, its been an hour and I'm already finished!!\r\nFather: oh sorry, by we, I meant you.",
			"thumbs_down": 60
		  }
		]
	}`
	failedBody := `{[]}`
	Client = &MockClient{}
	testcases := []struct {
		name               string
		term               string
		doFunc             mockDoFunc
		httpStatusExpected int
		lengthResponse     int
	}{
		{
			name: "success - definition by term",
			term: "ha",
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(successBody))),
				}, nil
			},
			httpStatusExpected: http.StatusOK,
			lengthResponse:     1,
		},
		{
			name: "failure - definition not found",
			term: "carnival",
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(failedBody))),
				}, nil
			},
			httpStatusExpected: http.StatusNotFound,
			lengthResponse:     0,
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			var list model.List
			Client = &MockClient{
				DoFunc: test.doFunc,
			}
			request, err := http.NewRequest(http.MethodGet, "www.arandomurl.com", nil)
			assert.Nil(t, err)
			response, err := Client.Do(request)
			assert.Nil(t, err)
			assert.Equal(t, test.httpStatusExpected, response.StatusCode)
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			assert.Nil(t, err)
			json.Unmarshal(body, &list)
			assert.Equal(t, test.lengthResponse, len(list.Definitions))

		})
	}
}

func TestGetDefinitionById(t *testing.T) {
	urbanClient := NewUrbanDictionary(testApiKey)
	testcases := []struct {
		name           string
		path           string
		id             string
		lengthResponse int
		error          error
		assertErr      func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertList     func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:           "success - definition found in CSV",
			path:           "../../data/definitions_test.csv",
			id:             "10593002",
			lengthResponse: 1,
			error:          nil,
			assertErr:      assert.Nil,
			assertList:     assert.NotNil,
		},
		{
			name:           "failure - definition not found in CSV",
			path:           "../../data/definitions_test.csv",
			id:             "1212",
			lengthResponse: 0,
			error:          model.ErrNotFoundInCSV{Id: "1212"},
			assertErr:      assert.NotNil,
			assertList:     assert.Nil,
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			urbanClient.UpdateCSVPath(test.path)
			definitions, err := urbanClient.GetDefinitionById(test.id)
			test.assertErr(t, err)
			test.assertList(t, definitions)
		})
	}
}

func TestOpen(t *testing.T) {
	// TODO add failure case
	urbanClient := NewUrbanDictionary(testApiKey)
	testcases := []struct {
		name      string
		path      string
		assertErr func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:      "success - open CSV file",
			path:      "../../data/definitions_test.csv",
			assertErr: assert.Nil,
		},
		{
			name:      "failure - could not open CSV file",
			path:      "../../data/definitions_fake.csv",
			assertErr: assert.NotNil,
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			urbanClient.UpdateCSVPath(test.path)
			_, err := urbanClient.Open()
			test.assertErr(t, err)
		})
	}
}

func TestRead(t *testing.T) {
	urbanClient := NewUrbanDictionary(testApiKey)
	testcases := []struct {
		name           string
		id             string
		path           string
		lengthResponse int
		assertErr      func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertLength   func(t assert.TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:           "success - reading content of the CSV",
			id:             "10593002",
			path:           "../../data/definitions_test.csv",
			lengthResponse: 1,
			assertErr:      assert.Nil,
			assertLength:   assert.Equal,
		},
		{
			name:           "failure - cannot read the content of the CSV",
			id:             "10593002",
			path:           "../../data/definitions_failure.csv",
			lengthResponse: 0,
			assertErr:      assert.NotNil,
			assertLength:   assert.Equal,
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			urbanClient.UpdateCSVPath(test.path)
			definitions, err := urbanClient.Read(test.id)
			test.assertErr(t, err)
			fmt.Println(len(definitions))
			test.assertLength(t, test.lengthResponse, len(definitions))
		})
	}
}

func TestWrite(t *testing.T) {
	urbanClient := NewUrbanDictionary(testApiKey)
	testcases := []struct {
		name              string
		path              string
		definitionToWrite *model.List
		assertErr         func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name: "success - reading content of the CSV",
			path: "../../data/definitions_test.csv",
			definitionToWrite: &model.List{
				Definitions: []model.Definition{
					{Defid: 12345, Word: "Hola"},
				},
			},
			assertErr: assert.Nil,
		},
		{
			name: "failure - cannot read the content of the CSV",
			path: "../../data/definitions_fake.csv",
			definitionToWrite: &model.List{
				Definitions: []model.Definition{
					{Defid: 89312, Word: "write"},
				},
			},
			assertErr: assert.NotNil,
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			urbanClient.UpdateCSVPath(test.path)
			err := urbanClient.Write(test.definitionToWrite)
			test.assertErr(t, err)

		})
	}
}
