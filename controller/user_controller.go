package controller

import (
	"net/http"

	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/model"
	"github.com/gkontos/goapi/util"
)

func (api *apiController) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := api.dbh.GetUsers()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error getting users")
		util.ReturnErrorJSON(w, err)
		return
	}
	if users == nil {
		users = make([]model.User, 0)
	}

	util.ReturnBodyJSON(w, users, http.StatusOK)

}
