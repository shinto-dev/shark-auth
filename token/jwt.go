package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"shark-auth/autherrors"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateJwtToken(userName string) (string, error) {
	expireAt := time.Now().Add(5 * time.Minute)
	claims := Claims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJwtToken(token string) (Claims, error) {
	var claims Claims
	tkn, err := jwt.ParseWithClaims(token, &claims, func(tkn *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logrus.Error("parsing token failed")
			return Claims{}, autherrors.ErrAuthenticationFailed
		}
		return Claims{}, autherrors.ErrInvalidToken
	}

	if !tkn.Valid {
		return Claims{}, autherrors.ErrAuthenticationFailed
	}

	return claims, nil
}
