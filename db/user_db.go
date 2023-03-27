package db

import (
	"database/sql"

	"github.com/gkontos/goapi/model"
	"github.com/google/uuid"
)

func (db *dbHandler) UpsertUser(u *model.User) (*model.User, error) {

	sqlStatement := ""
	if u.UID == "" {
		u.UID = uuid.NewString()
		sqlStatement = `
			INSERT INTO users (uid, auth_provider, provider_id, user_name, details, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW())`
	} else {
		sqlStatement = `
		UPDATE users
		SET auth_provider = $2, provider_id = $3, user_name = $4, details = $5, updated_at = NOW()
		WHERE uid = $1`
	}
	_, err := db.getConnection().Exec(sqlStatement, u.UID, u.AuthProvider, u.ProviderID, u.UserName, u.UserDetails)
	if err != nil {
		return nil, err
	} else {
		return u, nil
	}
}

func (db *dbHandler) GetUserByProvider(authProvider string, providerID string) (*model.User, error) {

	u := model.User{}
	sqlStatement := `
		SELECT uid, auth_provider, provider_id, user_name, details FROM users
		WHERE auth_provider = $1 AND provider_id = $2`
	err := db.getConnection().QueryRow(sqlStatement,
		authProvider,
		providerID).
		Scan(&u.UID,
			&u.AuthProvider,
			&u.ProviderID,
			&u.UserName,
			&u.UserDetails)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &u, nil
}

func (db *dbHandler) GetUsers() ([]model.User, error) {
	userActivity := make([]model.User, 0)
	sqlStatement := `
		SELECT uid, user_name, created_at, updated_at, details FROM users u 
		`
	rows, err := db.getConnection().Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.UID,
			&u.UserName,
			&u.CreatedDate,
			&u.LastLogin,
			&u.UserDetails,
		); err != nil {
			return userActivity, err
		}

		userActivity = append(userActivity, u)
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return userActivity, nil
}
