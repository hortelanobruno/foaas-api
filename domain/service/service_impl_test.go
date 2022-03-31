package service

import (
	"fmt"
	"github.com/hortelanobruno/foaas-api/domain/model"
	httpmock "github.com/hortelanobruno/foaas-api/http/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMessage(t *testing.T) {
	cases := []struct {
		name             string
		userID           string
		mockClient       *httpmock.Client
		expectedResponse *model.Response
		expectedError    error
	}{
		{
			"Should return an error when client returns an error",
			"123",
			func() *httpmock.Client {
				mock := &httpmock.Client{}
				mock.On("Get", "https://foaas.com/asshole/123").
					Return(nil, fmt.Errorf("error getting response from foaas"))
				return mock
			}(),
			nil,
			fmt.Errorf("error getting response from foaas"),
		},
		{
			"Should return an error when there's an error unmarshalling the response",
			"123",
			func() *httpmock.Client {
				mock := &httpmock.Client{}
				mock.On("Get", "https://foaas.com/asshole/123").
					Return(nil, nil)
				return mock
			}(),
			nil,
			fmt.Errorf("error unmarshaling the body, err: unexpected end of JSON input"),
		},
		{
			"Should return an error when there's an error unmarshalling the response",
			"123",
			func() *httpmock.Client {
				mock := &httpmock.Client{}
				mock.On("Get", "https://foaas.com/asshole/123").
					Return([]byte(`{"message": "Fuck you, asshole.","subtitle": "- 123"}`),
						nil)
				return mock
			}(),
			&model.Response{Message: "Fuck you, asshole.",
				Subtitle: "- 123"},
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			service := NewMessageServiceImpl(c.mockClient)

			// Operation
			response, err := service.GetMessage(c.userID)

			// Validation
			assert.EqualValues(t, c.expectedResponse, response)
			assert.EqualValues(t, c.expectedError, err)
			c.mockClient.AssertNumberOfCalls(t, "Get", 1)
		})
	}
}
