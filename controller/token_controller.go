package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/model"
	"github.com/gkontos/goapi/util"
)

type AuthenticationError struct {
	Err error
}

type TokenRequest struct {
	Token string `json:"token"`
}

func (e *AuthenticationError) Error() string {
	return e.Err.Error()
}

func (api *apiController) TokenCreate(w http.ResponseWriter, r *http.Request) {

	loginRequestToken := &TokenRequest{}
	if parseErr := util.ParseJsonRequest(r, &loginRequestToken); parseErr != nil {
		util.ReturnErrorJSON(w, parseErr)
		return
	}

	token, err := api.th.ValidateLoginAndCreateAccessToken(loginRequestToken.Token)

	if err != nil {
		logger.Logger.Error().Msg(fmt.Sprintf("unable to get auth token : %v", err))
		loginErr := &AuthenticationError{
			Err: errors.New("unable to process login request"),
		}
		util.ReturnErrorJSONWithCode(w, loginErr, http.StatusForbidden)

		return
	}

	util.ReturnBodyJSON(w, token, http.StatusOK)

}

func (api *apiController) TokenRefresh(w http.ResponseWriter, r *http.Request) {

	assertedUser := &model.Token{}
	if parseErr := util.ParseJsonRequest(r, &assertedUser); parseErr != nil {
		util.ReturnErrorJSON(w, parseErr)
		return
	}
	token, err := api.th.RefreshToken(assertedUser.RefreshToken)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("unable to get refresh token")
		loginErr := &AuthenticationError{
			Err: errors.New("unable to process refresh request"),
		}
		util.ReturnErrorJSONWithCode(w, loginErr, http.StatusForbidden)

		return
	}

	logger.Logger.Debug().Msg("new token issued.")
	util.ReturnBodyJSON(w, token, http.StatusOK)

}
