package http

import (
	"fmt"
	"github.com/hortelanobruno/foaas-api/constants"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type ClientImpl struct {
	client *http.Client
}

func NewClientImpl(timeout time.Duration) *ClientImpl {
	client := &http.Client{
		Timeout: timeout,
	}
	return &ClientImpl{
		client: client,
	}
}

func (c *ClientImpl) Get(url string) ([]byte, error) {
	logrus.Debugf("Starting to get response for %s", url)
	httpReq, err := c.generateHTTPRequest(url)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		logrus.Errorf("Error executing GET request for url: %s, err: %s", url, err.Error())
		return nil, fmt.Errorf("error doing the request, err: %s", err.Error())
	}

	body, err := c.processResponse(httpResp)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Finishing to get response %+v with body: %s for url: %s", httpResp,
		string(body), url)
	return body, nil
}

func (c *ClientImpl) GetWithUserIdHeader(url, userId string) ([]byte, error) {
	logrus.Debugf("Starting to get response for %s", url)
	httpReq, err := c.generateHTTPRequest(url)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(constants.UserIDHeader, userId)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		logrus.Errorf("Error executing GET request for url: %s, err: %s", url, err.Error())
		return nil, fmt.Errorf("error doing the request, err: %s", err.Error())
	}

	body, err := c.processResponse(httpResp)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Finishing to get response %+v with body: %s for url: %s", httpResp,
		string(body), url)
	return body, nil
}

func (c *ClientImpl) generateHTTPRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("Error creating request for url: %s, err: %s", url, err.Error())
		return nil, fmt.Errorf("error building the request, err: %s", err.Error())
	}
	req.Header.Add("Accept", "application/json")
	return req, nil
}

func (c *ClientImpl) processResponse(httpResp *http.Response) ([]byte, error) {
	body, err := c.readBody(httpResp)
	if err != nil {
		return nil, err
	}

	if http.StatusOK != httpResp.StatusCode {
		logrus.Errorf("Status code (%d) is different than OK", httpResp.StatusCode)
		return nil, fmt.Errorf("error executing request, status code: %v", httpResp.StatusCode)
	}

	return body, nil
}

func (c *ClientImpl) readBody(httpResp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		logrus.Errorf("Error reading body %v, err: %s", httpResp.Body, err.Error())
		return nil, fmt.Errorf("error reading the body, err: %s", err.Error())
	}

	err = httpResp.Body.Close()
	if err != nil {
		logrus.Errorf("Error closing body, err: %s", err.Error())
	}
	return body, nil
}
