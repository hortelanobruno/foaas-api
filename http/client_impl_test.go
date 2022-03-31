package http

import (
	"bytes"
	"fmt"
	httpmock "github.com/hortelanobruno/foaas-api/http/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReadBody(t *testing.T) {
	cases := []struct {
		name          string
		input         *http.Response
		expectedBody  []byte
		expectedError error
	}{
		{
			"Should return an error when readAll method returns an error",
			func() *http.Response {
				readCloserMock := &httpmock.ReadCloser{}
				readCloserMock.On("Read", mock.Anything).
					Return(0, fmt.Errorf("error reading"))
				httpResp := &http.Response{
					Body: readCloserMock,
				}
				return httpResp
			}(),
			nil,
			fmt.Errorf("error reading the body, err: error reading"),
		},
		{
			"Should return a nil error and the body",
			&http.Response{
				Body: ioutil.NopCloser(bytes.NewReader([]byte("unit testing"))),
			},
			[]byte("unit testing"),
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			client := NewClientImpl(time.Duration(0))

			// Operation
			body, err := client.readBody(c.input)

			// Validation
			assert.EqualValues(t, c.expectedBody, body)
			assert.EqualValues(t, c.expectedError, err)
		})
	}
}

func TestProcessResponse(t *testing.T) {
	cases := []struct {
		name          string
		input         *http.Response
		expectedBody  []byte
		expectedError error
	}{
		{
			"Should return an error when readAll method returns an error",
			func() *http.Response {
				readCloserMock := &httpmock.ReadCloser{}
				readCloserMock.On("Read", mock.Anything).
					Return(0, fmt.Errorf("error reading"))
				httpResp := &http.Response{
					Body: readCloserMock,
				}
				return httpResp
			}(),
			nil,
			fmt.Errorf("error reading the body, err: error reading"),
		},
		{
			"Should return an error when status code is different than 200",
			&http.Response{
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("unit testing"))),
				StatusCode: http.StatusInternalServerError,
			},
			nil,
			fmt.Errorf("error executing request, status code: 500"),
		},
		{
			"Should return a nil error and the body",
			&http.Response{
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("unit testing"))),
				StatusCode: http.StatusOK,
			},
			[]byte("unit testing"),
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			client := NewClientImpl(time.Duration(0))

			// Operation
			body, err := client.processResponse(c.input)

			// Validation
			assert.EqualValues(t, c.expectedBody, body)
			assert.EqualValues(t, c.expectedError, err)
		})
	}
}

func TestGenerateHTTPRequest(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expectedReq   *http.Request
		expectedError error
	}{
		{
			"Should return an error when url is invalid",
			":",
			nil,
			fmt.Errorf(`error building the request, err: parse ":": missing protocol scheme`),
		},
		{
			"Should return a nil error",
			"https://foaas.com/version",
			func() *http.Request {
				req, _ := http.NewRequest("GET", "https://foaas.com/version", nil)
				req.Header.Set("Accept", "application/json")
				return req
			}(),
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			client := NewClientImpl(time.Duration(0))

			// Operation
			req, err := client.generateHTTPRequest(c.input)

			// Validation
			assert.EqualValues(t, c.expectedReq, req)
			assert.EqualValues(t, c.expectedError, err)
		})
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RawQuery {
		case "status=ok":
			w.WriteHeader(http.StatusOK)
		case "status=badRequest":
			w.WriteHeader(http.StatusBadRequest)
		}
		_, _ = fmt.Fprint(w, `{"message": "ok"}`)
	}))
	defer server.Close()

	cases := []struct {
		name          string
		url           string
		expectedBody  []byte
		expectedError error
	}{
		{
			"Should return an error when url is invalid",
			":",
			nil,
			fmt.Errorf(`error building the request, err: parse ":": missing protocol scheme`),
		},
		{
			"Should return a nil error",
			fmt.Sprintf("%s?status=ok", server.URL),
			[]byte(`{"message": "ok"}`),
			nil,
		},
		{
			"Should return a nil error",
			fmt.Sprintf("%s?status=badRequest", server.URL),
			nil,
			fmt.Errorf("error executing request, status code: 400"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			client := NewClientImpl(time.Duration(0))

			// Operation
			body, err := client.Get(c.url)

			// Validation
			assert.EqualValues(t, c.expectedBody, body)
			assert.EqualValues(t, c.expectedError, err)
		})
	}
}

func TestGetWithUserIdHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.RawQuery {
		case "status=ok":
			w.WriteHeader(http.StatusOK)
		case "status=badRequest":
			w.WriteHeader(http.StatusBadRequest)
		}
		_, _ = fmt.Fprint(w, `{"message": "ok"}`)
	}))
	defer server.Close()

	cases := []struct {
		name          string
		url           string
		userID        string
		expectedBody  []byte
		expectedError error
	}{
		{
			"Should return an error when url is invalid",
			":",
			"123",
			nil,
			fmt.Errorf(`error building the request, err: parse ":": missing protocol scheme`),
		},
		{
			"Should return a nil error",
			fmt.Sprintf("%s?status=ok", server.URL),
			"123",

			[]byte(`{"message": "ok"}`),
			nil,
		},
		{
			"Should return a nil error",
			fmt.Sprintf("%s?status=badRequest", server.URL),
			"123",
			nil,
			fmt.Errorf("error executing request, status code: 400"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Initialization
			client := NewClientImpl(time.Duration(0))

			// Operation
			body, err := client.GetWithUserIdHeader(c.url, c.userID)

			// Validation
			assert.EqualValues(t, c.expectedBody, body)
			assert.EqualValues(t, c.expectedError, err)
		})
	}
}
