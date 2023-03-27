package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// --------------- security package -------------------- //
// User defines a user for our application for authentication, reference, and storage
type User struct {
	UID          string      `json:"uid"`
	AuthProvider string      `json:"auth_provider"`
	ProviderID   string      `json:"provider_id"`
	UserName     string      `json:"user_name"`
	CreatedDate  time.Time   `json:"created_date"`
	LastLogin    time.Time   `json:"last_login_date"`
	UserDetails  UserDetails `json:"user_details,omitempty"`
}

type UserDetails struct {
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Email     string   `json:"email,omitempty"`
	FullName  string   `json:"full_name"`
	Roles     []string `json:"roles"`
}

// HasRole returns true if the user is in the role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.UserDetails.Roles {
		if role == roleName {
			return true
		}
	}
	return false
}

// for pg jsonb value : Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (ud UserDetails) Value() (driver.Value, error) {
	return json.Marshal(ud)
}

// for pg jsonb value :Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (ud *UserDetails) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &ud)
}
