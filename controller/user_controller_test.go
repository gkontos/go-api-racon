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
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	logger.InitLogger(true, true)
	path := "/v1/users"

	userListResponse := []model.User{
		{
			UID:          "someuid",
			AuthProvider: "myspace",
			ProviderID:   "tom",
			UserName:     "butler",
			LastLogin:    time.Now(),
			UserDetails: model.UserDetails{
				Roles: []string{"ROLE_USER"},
			},
		},
	}
	cases := []struct {
		expectedResponseCode int
		expectedResponseBody []byte
		expectedResponseObj  []model.User
		ctrl                 *apiController
		requestBody          []byte
	}{
		// empty set
		{
			ctrl: &apiController{
				dbh: &testDbHandler{
					throwDbError: false,
					userResponse: nil,
				},
			},
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: []byte(`[]`),
		},
		// ok
		{
			ctrl: &apiController{
				dbh: &testDbHandler{
					throwDbError: false,
					userResponse: userListResponse,
				},
			},
			expectedResponseObj:  userListResponse,
			expectedResponseCode: http.StatusOK,
		},
		// db error
		{
			ctrl: &apiController{
				dbh: &testDbHandler{
					throwDbError: true,
					userResponse: userListResponse,
				},
			},
			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: []byte(`{"error":"db error"}`),
		},
	}

	for _, c := range cases {
		rootRequest, err := http.NewRequest("POST", path, bytes.NewBuffer(c.requestBody))
		if err != nil {
			t.Errorf("Root request error: %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.ctrl.GetUsers)
		handler.ServeHTTP(rr, rootRequest)

		assert.Equal(t, c.expectedResponseCode, rr.Code, "status code didn't match")

		logger.Logger.Info().Msg(rr.Body.String())
		if c.expectedResponseBody != nil {
			assert.Equal(t, string(c.expectedResponseBody), string(bytes.TrimSpace(rr.Body.Bytes())), "body didn't match")
		} else {
			var resp []model.User
			json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.True(t, len(resp) == 1)

		}
	}
}

type testDbHandler struct {
	userResponse []model.User
	throwDbError bool
}

// dbHandler implementation
func (d testDbHandler) GetUsers() ([]model.User, error) {
	if d.throwDbError {
		return nil, errors.New("db error")
	}
	return d.userResponse, nil
}

func (d testDbHandler) GetUserByProvider(authProvider string, providerID string) (*model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (d testDbHandler) UpsertUser(u *model.User) (*model.User, error) {
	panic("not implemented") // TODO: Implement
}
