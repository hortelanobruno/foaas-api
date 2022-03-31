package service

import (
	"encoding/json"
	"fmt"
	"github.com/hortelanobruno/foaas-api/constants"
	"github.com/hortelanobruno/foaas-api/domain/model"
	"github.com/hortelanobruno/foaas-api/http"
	"github.com/sirupsen/logrus"
)

type MessageServiceImpl struct {
	FoaasProtocol string
	FoaasDomain   string
	client        http.Client
}

func NewMessageServiceImpl(client http.Client) *MessageServiceImpl {
	return &MessageServiceImpl{
		FoaasProtocol: constants.FoaasProtocol,
		FoaasDomain:   constants.FoaasDomain,
		client:        client,
	}
}

func (m *MessageServiceImpl) GetMessage(userID string) (*model.Response, error) {
	url := fmt.Sprintf("%s://%s/asshole/%s", m.FoaasProtocol, m.FoaasDomain, userID)
	body, err := m.client.Get(url)
	if err != nil {
		return nil, err
	}

	response := &model.Response{}
	if err := json.Unmarshal(body, response); err != nil {
		logrus.Errorf("Error unmarshaling the response, err: %s", err.Error())
		return nil, fmt.Errorf("error unmarshaling the body, err: %s", err.Error())
	}

	return response, nil
}
