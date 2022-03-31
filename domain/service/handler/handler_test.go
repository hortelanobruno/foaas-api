package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hortelanobruno/foaas-api/domain/model"
	servicemocks "github.com/hortelanobruno/foaas-api/domain/service/mocks"
	validatormocks "github.com/hortelanobruno/foaas-api/domain/validator/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetMessage(t *testing.T) {
	cases := []struct {
		name                 string
		userID               string
		mockMessageValidator *validatormocks.MessageValidator
		mockMessageService   *servicemocks.MessageService
		expectedStatusCode   int
		expectedBody         string
	}{
		{
			"Should return an error when validator returns an error",
			"",
			func() *validatormocks.MessageValidator {
				mock := &validatormocks.MessageValidator{}
				mock.On("ValidateMessage", "").
					Return(fmt.Errorf("an error"))
				return mock
			}(),
			nil,
			http.StatusBadRequest,
			`{"error":"an error"}`,
		},
		{
			"Should return an error when service returns an error",
			"123",
			func() *validatormocks.MessageValidator {
				mock := &validatormocks.MessageValidator{}
				mock.On("ValidateMessage", "123").
					Return(nil)
				return mock
			}(),
			func() *servicemocks.MessageService {
				mock := &servicemocks.MessageService{}
				mock.On("GetMessage", "123").
					Return(nil, fmt.Errorf("error getting message"))
				return mock
			}(),
			http.StatusInternalServerError,
			`{"error":"error getting message"}`,
		},
		{
			"Should return a nil error",
			"123",
			func() *validatormocks.MessageValidator {
				mock := &validatormocks.MessageValidator{}
				mock.On("ValidateMessage", "123").
					Return(nil)
				return mock
			}(),
			func() *servicemocks.MessageService {
				mock := &servicemocks.MessageService{}
				mock.On("GetMessage", "123").
					Return(&model.Response{
						Message:  "message",
						Subtitle: "subtitle",
					}, nil)
				return mock
			}(),
			http.StatusOK,
			`{"message":"message","subtitle":"subtitle"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			handler := NewMessageHandler(c.mockMessageValidator, c.mockMessageService)

			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)
			context.Request, _ = http.NewRequest("GET", "/", nil)
			context.Request.Header.Set("UserId", c.userID)

			// Operation
			handler.HandleGetMessage(context)

			// Validation
			assert.EqualValues(t, c.expectedStatusCode, w.Code)
			assert.EqualValues(t, c.expectedBody, w.Body.String())

			if c.mockMessageValidator != nil {
				c.mockMessageValidator.AssertNumberOfCalls(t, "ValidateMessage", 1)
			}

			if c.mockMessageService != nil {
				c.mockMessageService.AssertNumberOfCalls(t, "GetMessage", 1)
			}
		})
	}
}
