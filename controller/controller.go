package controller

import (
	"github.com/gkontos/goapi/db"
	"github.com/gkontos/goapi/security"
	"github.com/go-chi/chi/v5"
)

type ApiRouter interface {
	SetupRouter() *chi.Mux
}
type apiRouter struct {
	rs   security.RouterSecurity
	ctrl *apiController
}

type apiController struct {
	dbh db.DbHandler
	th  security.TokenHandler
}

func NewRouter(allow_origin string, apiController *apiController, dbHandler db.DbHandler) ApiRouter {
	return &apiRouter{
		rs:   security.NewRouterSecurity(allow_origin, dbHandler),
		ctrl: apiController,
	}
}

func NewController(dbHandler db.DbHandler) *apiController {
	return &apiController{
		dbh: dbHandler,
		th:  security.GetNewHandler(dbHandler),
	}
}
