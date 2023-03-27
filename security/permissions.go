package security

import (
	"errors"

	"github.com/gkontos/goapi/model"
)

type Permission string

func (p Permission) String() string {
	return string(p)
}

func (s *tokenHandler) CheckPermission(user *model.User, permission Permission) error {
	if user == nil {
		return errors.New("CheckPermission: No user supplied")
	}

	if permission == "" {
		return errors.New("CheckPermission: You must supply a valid permission to check against")
	}

	if user.HasRole(AdministratorRole) {
		// Admins can do anything
		return nil
	}

	if permission.String() == "read" || permission.String() == "write" && user.HasRole(UserRole) {
		// User has role
		return nil
	}

	return errors.New("CheckPermission: User not authorized")
}
