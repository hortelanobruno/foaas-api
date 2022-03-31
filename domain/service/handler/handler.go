package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hortelanobruno/foaas-api/constants"
	"github.com/hortelanobruno/foaas-api/domain/service"
	"github.com/hortelanobruno/foaas-api/domain/validator"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MessageHandler struct {
	messageValidator validator.MessageValidator
	messageService   service.MessageService
}

func NewMessageHandler(messageValidator validator.MessageValidator,
	messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageValidator: messageValidator,
		messageService:   messageService,
	}
}

func (m *MessageHandler) HandleGetMessage(ginContext *gin.Context) {
	userID := ginContext.GetHeader(constants.UserIDHeader)
	if err := m.messageValidator.ValidateMessage(userID); err != nil {
		logrus.Errorf("Error validating the message, userID: %s", userID)
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := m.messageService.GetMessage(userID)
	if err != nil {
		logrus.Errorf("Error getting the message, userID: %s, err: %s", userID, err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ginContext.JSON(http.StatusOK, response)
	return
}
