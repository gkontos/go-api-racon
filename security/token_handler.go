package security

import (
	"errors"
	"time"

	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/model"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UID       string   `json:"uid"`
	Username  string   `json:"user_name"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	Activated bool     `json:"activated"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Image     string   `json:"image"`
	jwt.RegisteredClaims
}

// ValidateLoginAndCreateToken will
// validate token; add / update user; create access tokens for the local issuer
// assumes the tokenrequest is a google auth token
func (s *tokenHandler) ValidateLoginAndCreateAccessToken(t string) (*model.Token, error) {

	googleClaims, err := validateGoogleJWT(t)
	if err != nil {
		return nil, err
	}
	// SAVE / UPDATE USER
	claims := mapGoogleClaimToClaims(googleClaims)

	if claims.Issuer == "" || claims.ID == "" {
		return nil, errors.New("issuer and ID are required claim fields")
	}

	user, err := s.createOrUpdateLocalUser(claims)
	if err != nil {
		return nil, err
	}

	claims.UID = user.UID

	// get tokens
	return s.obtainAccessTokens(claims)

}

func (s *tokenHandler) createOrUpdateLocalUser(claims Claims) (*model.User, error) {

	user, err := s.dbh.GetUserByProvider(claims.Issuer, claims.ID)
	if err != nil {
		return nil, err
	}

	claims.Roles = user.UserDetails.Roles

	u := CreateUserFromClaims(claims)
	if user.UID != "" {
		u.UID = user.UID
	} else {
		// add a default role for new users
		u.UserDetails.Roles = append(u.UserDetails.Roles, UserRole)
	}

	user, err = s.dbh.UpsertUser(u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *tokenHandler) RefreshToken(token string) (*model.Token, error) {

	claims, err := s.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}
	return s.obtainAccessTokens(claims)

}

func (s *tokenHandler) ValidateAccessToken(tokenString string) (Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure token's signature wasn't changed
		return verificationKey, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return *claims, nil
	}

	return Claims{}, errors.New("invalid token")

}

// create a local jwt token
func (s *tokenHandler) obtainAccessTokens(claims Claims) (*model.Token, error) {

	token_expires_at := time.Now().Add(time.Minute * time.Duration(tokenExpirationMinutes))
	refresh_expires_at := time.Now().Add(time.Minute * time.Duration(refreshTokenExpirationMinutes))
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(token_expires_at),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "local",
		Subject:   "access",
		Audience:  []string{"local"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("unable to create token")
		return nil, err
	}

	refreshClaims := claims
	refreshClaims.RegisteredClaims = claims.RegisteredClaims
	refreshClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(refresh_expires_at)

	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims).SignedString(signingKey)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("unable to create token")
		return nil, err
	}

	return &model.Token{
		Token:        signedToken,
		RefreshToken: refresh_token,
		ExpiresAt:    token_expires_at,
	}, nil

}
