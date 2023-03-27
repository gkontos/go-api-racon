package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/model"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DbHandler interface {
	GetUsers() ([]model.User, error)
	GetUserByProvider(authProvider string, providerID string) (*model.User, error)
	UpsertUser(u *model.User) (*model.User, error)
}

type dbHandler struct {
	dbUser         string
	dbPassword     string
	unixSocketPath string
	dbName         string
	pool           *sql.DB
}

/*
*
setup logger, get connection params from secrets/env
*
*/
func NewDbHandler() *dbHandler {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			logger.Logger.Error().Msg(fmt.Sprintf("Error: %s environment variable not set.", k))
			panic("env variables not set for db conx")
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.

	return &dbHandler{
		dbUser:         mustGetenv("DB_USER"),              // e.g. 'my-db-user'
		dbPassword:     mustGetenv("DB_PASS"),              // e.g. 'my-db-password'
		unixSocketPath: mustGetenv("INSTANCE_UNIX_SOCKET"), // e.g. '/cloudsql/project:region:instance'
		dbName:         mustGetenv("DB_NAME"),              // e.g. 'my-database'
	}
}

/*
*
get a sql connection
*
*/
func (db *dbHandler) getConnection() *sql.DB {
	var err error
	pool, err := db.connectUnixSocket()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Connection error")
		return nil
	}
	db.pool = pool
	return pool

}

// connectUnixSocket initializes a Unix socket connection pool for
// a Cloud SQL instance of Postgres.
func (db *dbHandler) connectUnixSocket() (*sql.DB, error) {
	if db.pool != nil {
		err := db.pool.Ping()
		if err == nil {
			return db.pool, nil
		}
	}

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		db.dbUser, db.dbPassword, db.dbName, db.unixSocketPath)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	err = dbPool.Ping()
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error establishing connection pool")
		return nil, err
	}

	return dbPool, nil
}
