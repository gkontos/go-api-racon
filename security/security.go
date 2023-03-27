package security

import (
	"crypto/rsa"
	"fmt"
	"os"
	"strconv"

	"github.com/gkontos/goapi/db"
	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/model"
	jwt "github.com/golang-jwt/jwt/v4"
)

var (
	// signingKey is the private key that will be used to sign jwt tokens on creation
	signingKey *rsa.PrivateKey
	// verificationKey is the public key that will be used to validate jwt tokens
	verificationKey               *rsa.PublicKey
	googleTokenAudience           string
	googleCerts                   map[string]*rsa.PublicKey
	tokenExpirationMinutes        int
	refreshTokenExpirationMinutes int
	tokenGraceSeconds             int
)

const (
	UserContextKey    = "usercxt"
	AdministratorRole = "ROLE_ADMIN"
	UserRole          = "ROLE_USER"
)

type TokenHandler interface {
	ValidateLoginAndCreateAccessToken(t string) (*model.Token, error)
	RefreshToken(t string) (*model.Token, error)
	ValidateAccessToken(tokenString string) (Claims, error)
}
type tokenHandler struct {
	dbh db.DbHandler
}

func initModule() {
	if verificationKey == nil {
		verificationKey = getRsaVerificationKey()
	}
	if signingKey == nil {
		signingKey = getRsaSigningKey()
	}

	if googleTokenAudience == "" {
		googleTokenAudience = mustGetenv("GOOGLE_TOKEN_AUDIENCE")
	}
	if tokenExpirationMinutes == 0 {
		tokenExpirationMinutes = getenvOrInt("TOKEN_VALID_MINUTES", 5)
	}
	if refreshTokenExpirationMinutes == 0 {
		refreshTokenExpirationMinutes = getenvOrInt("REFRESH_TOKEN_VALID_MINUTES", 10)
	}
	if tokenGraceSeconds == 0 {
		tokenGraceSeconds = getenvOrInt("TOKEN_GRACE_SECONDS", 5)
	}
	if googleCerts == nil {
		googleCerts = make(map[string]*rsa.PublicKey)
	}
}
func GetNewHandler(dbHandler db.DbHandler) *tokenHandler {
	initModule()
	return &tokenHandler{
		dbh: dbHandler,
	}
}

// TODO: this should be coming from a secure source and not env variables
func getRsaSigningKey() *rsa.PrivateKey {
	signBytes := []byte(mustGetenv("PRIVATE_KEY"))
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		logger.Logger.Error().Msg(fmt.Sprintf("error getting key from PEM %v", err))
		return nil
	}
	return signKey
}

// TODO: this should be coming from a secure source and not env variables
func getRsaVerificationKey() *rsa.PublicKey {
	verifyBytes := []byte(mustGetenv("PUBLIC_KEY"))
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		logger.Logger.Error().Msg(fmt.Sprintf("error getting key from PEM %v", err))
		return nil
	}
	return verifyKey
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		logger.Logger.Error().Msg(fmt.Sprintf("Warning: %s environment variable not set.\n", k))
		panic("env variables not set for tokens")
	}
	return v
}

func getenvOrInt(k string, defaultValue int) int {
	v := os.Getenv(k)
	if v == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(v)
	if err != nil {
		logger.Logger.Error().Err(err)
		return defaultValue
	}
	return val
}

func CreateUserFromClaims(claims Claims) *model.User {
	user := &model.User{
		UID:          claims.UID,
		AuthProvider: claims.Issuer,
		ProviderID:   claims.ID,
		UserName:     claims.Username,
		UserDetails: model.UserDetails{
			FirstName: claims.FirstName,
			LastName:  claims.LastName,
			Email:     claims.Email,
			Roles:     claims.Roles,
		},
	}
	return user
}
