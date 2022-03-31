package service

import (
	"github.com/hortelanobruno/foaas-api/domain/model"
)

type MessageService interface {
	GetMessage(userID string) (*model.Response, error)
}
