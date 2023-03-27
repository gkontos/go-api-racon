package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/model"
	"github.com/gkontos/goapi/security"
	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	logger.InitLogger(true, true)
	path := "/v1/login"

	cases := []struct {
		expectedResponseCode int
		expectedResponseBody []byte
		ctrl                 *apiController
		requestBody          []byte
	}{
		// bad format
		{
			ctrl: &apiController{
				th: &testTokenHandler{returnError: false},
			},
			requestBody:          []byte(`{"badformat"}`),
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: []byte(`{"error":"unable to parse request invalid character '}' after object key"}`),
		},
		// ok
		{
			ctrl: &apiController{
				th: &testTokenHandler{returnError: false},
			},
			requestBody:          []byte(`{"token":"sometoken"}`),
			expectedResponseCode: http.StatusOK,
		},
		// fail
		{
			ctrl: &apiController{
				th: &testTokenHandler{returnError: true},
			},
			requestBody:          []byte(`{"token": "sometoken"}`),
			expectedResponseCode: http.StatusForbidden,
			expectedResponseBody: []byte(`{"error":"unable to process login request"}`),
		},
	}

	for _, c := range cases {
		rootRequest, err := http.NewRequest("POST", path, bytes.NewBuffer(c.requestBody))
		if err != nil {
			t.Errorf("Root request error: %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.ctrl.TokenCreate)
		handler.ServeHTTP(rr, rootRequest)

		assert.Equal(t, c.expectedResponseCode, rr.Code, "status code didn't match")

		logger.Logger.Info().Msg(rr.Body.String())
		if c.expectedResponseBody != nil {
			assert.Equal(t, string(c.expectedResponseBody), string(bytes.TrimSpace(rr.Body.Bytes())), "body didn't match")
		} else {
			var resp model.Token
			json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.True(t, resp.RefreshToken != "")
			assert.True(t, resp.Token != "")
			assert.True(t, resp.ExpiresAt.String() != "")

		}
	}
}

func TestRefreshToken(t *testing.T) {
	logger.InitLogger(true, true)
	path := "/v1/login/refresh"

	cases := []struct {
		expectedResponseCode int
		expectedResponseBody []byte
		ctrl                 *apiController
		requestBody          []byte
	}{
		// bad format
		{
			ctrl: &apiController{
				th: &testTokenHandler{returnError: false},
			},
			requestBody:          []byte(`{"badformat"}`),
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: []byte(`{"error":"unable to parse request invalid character '}' after object key"}`),
		},
		// ok
		{
			ctrl: &apiController{
				th: &testTokenHandler{returnError: false},
			},
			requestBody:          []byte(`{"token":"sometoken", "refresh_token":"somerefreshvalue", "expires_at":"2006-03-17T15:04:05Z"}`),
			expectedResponseCode: http.StatusOK,
		},
		// fail
		{
			ctrl: &apiController{
				th: &testTokenHandler{returnError: true},
			},
			requestBody:          []byte(`{"token":"someinvalidtoken", "refresh_token":"someinvalidrefreshvalue", "expires_at":"2006-03-17T15:04:05Z"}`),
			expectedResponseCode: http.StatusForbidden,
			expectedResponseBody: []byte(`{"error":"unable to process refresh request"}`),
		},
	}

	for _, c := range cases {
		rootRequest, err := http.NewRequest("POST", path, bytes.NewBuffer(c.requestBody))
		if err != nil {
			t.Errorf("Root request error: %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.ctrl.TokenRefresh)
		handler.ServeHTTP(rr, rootRequest)

		assert.Equal(t, c.expectedResponseCode, rr.Code, "status code didn't match")

		logger.Logger.Info().Msg(rr.Body.String())
		if c.expectedResponseBody != nil {
			assert.Equal(t, string(c.expectedResponseBody), string(bytes.TrimSpace(rr.Body.Bytes())), "body didn't match")
		} else {
			var resp model.Token
			json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.True(t, resp.RefreshToken != "")
			assert.True(t, resp.Token != "")
			assert.True(t, resp.ExpiresAt.String() != "")

		}
	}
}

type testTokenHandler struct {
	returnError bool
}

// tokenHandler implementation
func (h *testTokenHandler) ValidateLoginAndCreateAccessToken(t string) (*model.Token, error) {
	if h.returnError {
		return nil, errors.New("some error")
	}
	return &model.Token{
		Token:        "sometokenstring",
		RefreshToken: "somerefreshtokenstring",
		ExpiresAt:    time.Now(),
	}, nil
}

func (h *testTokenHandler) RefreshToken(t string) (*model.Token, error) {
	if h.returnError {
		return nil, errors.New("refresh error")
	}
	return &model.Token{
		Token:        "sometokenstring",
		RefreshToken: "somerefreshtokenstring",
		ExpiresAt:    time.Now(),
	}, nil
}

func (h *testTokenHandler) ValidateAccessToken(tokenString string) (security.Claims, error) {
	return security.Claims{}, nil
}
