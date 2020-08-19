package token

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/sirupsen/logrus"

	"shark-auth/autherrors"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username     string `json:"username"`
	RefreshToken string `json:"refresh_uuid,omitempty"`
	jwt.StandardClaims
}

func CreateAccessToken(userName string) (string, error) {
	expireAt := time.Now().Add(5 * time.Minute)
	claims := Claims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expireAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseAccessToken(token string) (Claims, error) {
	var claims Claims
	tkn, err := jwt.ParseWithClaims(token, &claims, func(tkn *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		logrus.Errorf("parsing token failed: %v", err)
		if err == jwt.ErrSignatureInvalid {
			return Claims{}, autherrors.ErrAuthenticationFailed
		}
		// todo this includes token expired error
		return Claims{}, autherrors.ErrInvalidToken
	}

	if !tkn.Valid {
		return Claims{}, autherrors.ErrAuthenticationFailed
	}

	isSignedout, err := checkUserIsAlreadySignedOut(token)
	if err != nil {
		return Claims{}, err
	}
	if isSignedout {
		return Claims{}, autherrors.ErrAuthenticationFailed
	}

	return claims, nil
}

func checkUserIsAlreadySignedOut(token string) (bool, error) {
	isBlacklisted, err := IsAccessTokenBlacklisted(token)
	if err != nil {
		// todo also add the cause
		return true, autherrors.ErrInternal
	}

	if isBlacklisted {
		return true, nil
	}
	return false, nil
}

func DeleteAccessToken(accessToken string) error {
	claims, err := ParseAccessToken(accessToken)
	if err != nil {
		return autherrors.ErrAuthenticationFailed
	}

	return BlacklistAccessToken(accessToken, claims.ExpiresAt.Time)
}
