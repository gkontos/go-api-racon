package security

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gkontos/goapi/db"
	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/model"
	"github.com/gkontos/goapi/util"
	"golang.org/x/net/context"
)

type RouterSecurity interface {
	CorsHeaders(next http.Handler) http.Handler
	AuthenticateAuthHeader(next http.Handler) http.Handler
	Authorize(permissions ...Permission) func(next http.HandlerFunc) http.HandlerFunc
}

type defaultRouterSecurity struct {
	allowedOrigin string
	dbh           db.DbHandler
}

func NewRouterSecurity(allow_origin string, dbHandler db.DbHandler) RouterSecurity {
	return &defaultRouterSecurity{
		allowedOrigin: allow_origin,
		dbh:           dbHandler,
	}
}

// SecureHeaders adds secure headers to the API
func (s *defaultRouterSecurity) CorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", s.allowedOrigin)

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		}

		next.ServeHTTP(w, r)
	})
}

func (s *defaultRouterSecurity) AuthenticateAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Get token from the Authorization header
		// format: Authorization: Bearer

		// TODO : case sensitivity in header
		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			tokenString = tokens[0]
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}
		// If the token is empty...
		if tokenString == "" {
			// If we get here, the required token is missing
			loginErr := &model.AuthenticationError{
				Err: errors.New(http.StatusText(http.StatusUnauthorized)),
			}
			util.ReturnErrorJSONWithCode(w, loginErr, http.StatusUnauthorized)
			return
		}

		s := GetNewHandler(s.dbh)
		claims, err := s.ValidateAccessToken(tokenString)

		if err != nil {
			logger.Logger.Error().Err(err).Msg("parse with claims error")
			loginErr := &model.AuthenticationError{
				Err: errors.New(http.StatusText(http.StatusUnauthorized)),
			}
			util.ReturnErrorJSONWithCode(w, loginErr, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

// Authorize provides authorization middleware for our handlers
func (s *defaultRouterSecurity) Authorize(permissions ...Permission) func(next http.HandlerFunc) http.HandlerFunc {
	return s.AuthorizeModel("", permissions...)
}
func (s *defaultRouterSecurity) AuthorizeModel(accessModel string, permissions ...Permission) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := r.Context().Value(UserContextKey).(Claims)
			user := CreateUserFromClaims(claims)
			if user.UserName == "" {
				loginErr := &model.AuthenticationError{
					Err: errors.New(http.StatusText(http.StatusForbidden)),
				}
				util.ReturnErrorJSONWithCode(w, loginErr, http.StatusUnauthorized)
				return
			}
			for _, permission := range permissions {
				s := GetNewHandler(s.dbh)
				if err := s.CheckPermission(user, permission); err != nil {
					logger.Logger.Error().Err(err).Msg("error reading permissions")
					loginErr := &model.AuthenticationError{
						Err: errors.New(http.StatusText(http.StatusForbidden)),
					}
					util.ReturnErrorJSONWithCode(w, loginErr, http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
