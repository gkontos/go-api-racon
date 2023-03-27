package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gkontos/goapi/model"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := NewRouter("http://localhost", &apiController{}, &routerTestDbHandler{}).SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

type routerTestDbHandler struct{}

func (d *routerTestDbHandler) GetUsers() ([]model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (d *routerTestDbHandler) GetUserByProvider(authProvider string, providerID string) (*model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (d *routerTestDbHandler) UpsertUser(u *model.User) (*model.User, error) {
	panic("not implemented") // TODO: Implement
}
