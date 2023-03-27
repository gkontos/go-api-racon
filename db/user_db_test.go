package db

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gkontos/goapi/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	db, mock := getTestHandler()
	defer db.pool.Close()

	rows := sqlmock.NewRows([]string{"uid", "user_name", "created_at", "updated_at", "details"}).
		AddRow(uuid.NewString(), "sergei", time.Now(), time.Now(), []byte(`{"email":"sergei"}`)).
		AddRow(uuid.NewString(), "tom", time.Now(), time.Now(), []byte(`{"email":"tom"}`))

	mock.ExpectQuery("SELECT (.+) FROM users (.+)").WithArgs().WillReturnRows(rows)

	var users []model.User
	var err error
	if users, err = db.GetUsers(); err != nil {
		t.Errorf("error '%s' was not expected, while getting users", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.True(t, len(users) == 2)
}

func TestGetUsersByProvider(t *testing.T) {
	db, mock := getTestHandler()
	defer db.pool.Close()

	authProvider := "google"
	providerId := "12345"
	rows := sqlmock.NewRows([]string{"uid", "auth_provider", "provider_id", "user_name", "details"}).
		AddRow(uuid.NewString(), "google", "12345", "sergei", []byte(`{"email":"sergei"}`)).
		AddRow(uuid.NewString(), "myspace", "6789", "tom", []byte(`{"email":"tom"}`))

	mock.ExpectQuery("SELECT (.+) FROM users (.+)").WithArgs(authProvider, providerId).WillReturnRows(rows)

	var user *model.User
	var err error
	if user, err = db.GetUserByProvider(authProvider, providerId); err != nil {
		t.Errorf("error '%s' was not expected, while getting users", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.True(t, user.UserName == "sergei")
}

func getTestHandler() (*dbHandler, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	dbh := &dbHandler{
		pool: db,
	}
	return dbh, mock
}
