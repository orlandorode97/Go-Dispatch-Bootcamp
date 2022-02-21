package presenter

import (
	"testing"

	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestResponseDefinitions(t *testing.T) {
	testcases := []struct {
		name           string
		message        string
		response       *model.List
		error          error
		assertErr      func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertResponse func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:    "success - written on field with correct format",
			message: "written on field parsed successfully",
			error:   nil,
			response: &model.List{
				Definitions: []model.Definition{
					{
						Word:      "the",
						WrittenOn: "2006-04-30T20:18:42.000Z",
					},
				},
			},
			assertErr:      assert.Nil,
			assertResponse: assert.NotNil,
		},
		{
			name:    "failure - written on field with incorrect format",
			message: "written on field parsed unsuccessfully",
			error:   model.ErrParsingDate{"2018-04-26T19", UrbanLayout},
			response: &model.List{
				Definitions: []model.Definition{
					{
						Word:      "the",
						WrittenOn: "2018-04-26T19",
					},
				},
			},
			assertErr:      assert.NotNil,
			assertResponse: assert.Empty,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			presenter := NewDefinitionPresenter()
			response, err := presenter.ResponseDefinitions(test.response)
			test.assertErr(t, err, test.message)
			test.assertResponse(t, response)
		})
	}
}
