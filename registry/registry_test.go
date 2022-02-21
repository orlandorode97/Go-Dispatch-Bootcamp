package registry

import (
	"testing"

	"github.com/orlandorode97/go-disptach/infrastructure/api"
	"github.com/orlandorode97/go-disptach/interface/controller"
	"github.com/stretchr/testify/assert"
)

var (
	testApiKey = "cc2464e4a1mshb5ceeca91e5a6adp1fa80bjsn4b48e2408b87"
)

func TestNewRegistry(t *testing.T) {
	testcases := []struct {
		name              string
		message           string
		urbanClient       *api.UrbanDictionary
		assertUrbanClient func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:              "success - urban client set up",
			message:           "urban client in registry set up",
			urbanClient:       api.NewUrbanDictionary(testApiKey),
			assertUrbanClient: assert.NotNil,
		},
		{
			name:              "failure - urban client nil",
			message:           "registry app with a nil urban client",
			urbanClient:       nil,
			assertUrbanClient: assert.Nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			test.assertUrbanClient(t, test.urbanClient, test.message)
		})
	}
}

func TestNewAppController(t *testing.T) {
	appController := NewRegistry(nil).NewAppController()
	testcases := []struct {
		name         string
		message      string
		typeExpected controller.AppController
		typeReceived interface{}
		assertType   func(t assert.TestingT, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:         "success - return AppController type",
			message:      "AppController type received",
			typeExpected: appController,
			typeReceived: appController,
			assertType:   assert.IsType,
		},
		{
			name:         "failure - return another type",
			message:      "expected AppController type",
			typeExpected: controller.AppController{},
			typeReceived: struct{}{},
			assertType:   assert.NotSame,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			test.assertType(t, test.typeExpected, test.typeReceived, test.message)
		})
	}
}
