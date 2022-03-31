package integration

import (
	"encoding/json"
	"fmt"
	"github.com/hortelanobruno/foaas-api/cmd/server"
	"github.com/hortelanobruno/foaas-api/domain/model"
	"github.com/hortelanobruno/foaas-api/domain/service"
	"github.com/hortelanobruno/foaas-api/domain/service/handler"
	"github.com/hortelanobruno/foaas-api/domain/validator"
	customhttp "github.com/hortelanobruno/foaas-api/http"
	"github.com/hortelanobruno/foaas-api/ratelimiter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestIntegrationShouldReturnNoErrorWhenRateLimitIsNotExceeded(t *testing.T) {
	// Initialization
	userID := "123"
	foaasServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"message": "Fuck you, asshole.","subtitle": "- %s"}`, userID)
	}))
	defer foaasServer.Close()

	rateLimiter := ratelimiter.NewLocalRateLimiter(2, time.Millisecond*time.Duration(10000))
	httpClient := customhttp.NewClientImpl(time.Duration(5) * time.Second)
	messageService := service.NewMessageServiceImpl(httpClient)
	messageService.FoaasProtocol = "http"
	messageService.FoaasDomain = strings.Split(foaasServer.URL, "//")[1]
	messageValidator := validator.NewMessageValidatorImpl()
	messageHandler := handler.NewMessageHandler(messageValidator, messageService)
	serverPort := 4000
	serverUrl := fmt.Sprintf("http://localhost:%d/message", serverPort)

	go func() {
		server := server.NewServer(messageHandler, rateLimiter)
		server.Start(serverPort)
	}()

	// Operation
	responseAttempt1, errorAttempt1 := requestMessageForUser(httpClient, serverUrl, userID)
	responseAttempt2, errorAttempt2 := requestMessageForUser(httpClient, serverUrl, userID)

	// Validation
	assertValidResponse(t, responseAttempt1, errorAttempt1)
	assertValidResponse(t, responseAttempt2, errorAttempt2)
}

func TestIntegrationShouldReturnErrorWhenRateLimitIsExceeded(t *testing.T) {
	// Initialization
	userID := "123"
	foaasServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `{"message": "Fuck you, asshole.","subtitle": "- %s"}`, userID)
	}))
	defer foaasServer.Close()

	rateLimiter := ratelimiter.NewLocalRateLimiter(2, time.Millisecond*time.Duration(10000))
	httpClient := customhttp.NewClientImpl(time.Duration(5) * time.Second)
	messageService := service.NewMessageServiceImpl(httpClient)
	messageService.FoaasProtocol = "http"
	messageService.FoaasDomain = strings.Split(foaasServer.URL, "//")[1]
	messageValidator := validator.NewMessageValidatorImpl()
	messageHandler := handler.NewMessageHandler(messageValidator, messageService)
	serverPort := 4001
	serverUrl := fmt.Sprintf("http://localhost:%d/message", serverPort)

	go func() {
		server := server.NewServer(messageHandler, rateLimiter)
		server.Start(serverPort)
	}()

	// Operation
	responseAttempt1, errorAttempt1 := requestMessageForUser(httpClient, serverUrl, userID)
	responseAttempt2, errorAttempt2 := requestMessageForUser(httpClient, serverUrl, userID)
	responseAttempt3, errorAttempt3 := requestMessageForUser(httpClient, serverUrl, userID)

	// Validation
	assertValidResponse(t, responseAttempt1, errorAttempt1)
	assertValidResponse(t, responseAttempt2, errorAttempt2)
	assertNotValidResponse(t, responseAttempt3, errorAttempt3)
}

func assertNotValidResponse(t *testing.T, response *model.Response, err error) {
	assert.NotNil(t, err)
	assert.Nil(t, response)
}

func assertValidResponse(t *testing.T, response *model.Response, err error) {
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, "Fuck you, asshole.", response.Message)
	assert.EqualValues(t, "- 123", response.Subtitle)
}

func requestMessageForUser(httpClient *customhttp.ClientImpl, serverUrl, userId string) (*model.Response, error) {
	body, err := httpClient.GetWithUserIdHeader(serverUrl, userId)
	if err != nil {
		return nil, err
	}
	response := &model.Response{}
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}
	return response, nil
}
