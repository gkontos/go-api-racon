package security

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gkontos/goapi/logger"
	"github.com/golang-jwt/jwt/v4"
)

type GoogleClaims struct {
	GID           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	FullName      string `json:"name"`
	Image         string `json:"picture"`
	jwt.RegisteredClaims
}

// TODO; make parse / validation more generic.
// The jwt library contains ParseUnverified
// which can be used to get the issuer w/o actually
// parsing the token.  That would allow an implementation based on provider to be used
func validateGoogleJWT(tokenString string) (GoogleClaims, error) {
	claimsStruct := GoogleClaims{}
	jwt.TimeFunc = func() time.Time {
		return time.Now().UTC().Add(time.Second * time.Duration(tokenGraceSeconds))
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {

			key, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))

			if err != nil {
				logger.Logger.Error().Msg(fmt.Sprintf("error parsing key w/ kid: %v -- %v", token.Header["kid"], err))
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error parsing claims")
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("invalid token")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if !claims.VerifyAudience(googleTokenAudience, true) {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if !claims.VerifyExpiresAt(time.Now(), true) {
		return GoogleClaims{}, errors.New("token is expired")
	}

	return *claims, nil
}

func mapGoogleClaimToClaims(claims GoogleClaims) Claims {

	return Claims{
		Username:  claims.FullName,
		Email:     claims.Email,
		Activated: claims.EmailVerified,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Image:     claims.Image,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:     claims.GID,
			Issuer: claims.Issuer,
		},
	}
}

func getGooglePublicKey(keyId string) (*rsa.PublicKey, error) {

	key, ok := googleCerts[keyId]
	if ok {
		return key, nil
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return nil, err
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return nil, err
	}
	pem, ok := myResp[keyId]
	if !ok {
		return nil, errors.New("key not found")
	}

	key, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	if err != nil {
		logger.Logger.Error().Msg(fmt.Sprintf("error parsing key w/ kid: %v -- %v", keyId, err))
		return nil, err
	}

	googleCerts[keyId] = key
	return key, nil
}
