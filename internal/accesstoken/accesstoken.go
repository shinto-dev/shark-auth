package accesstoken

import (
	"shark-auth/internal/apperror"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/sirupsen/logrus"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"sid"`
	jwt.StandardClaims
}

func Create(userID string, sessionID string) (string, error) {
	expireAt := time.Now().Add(5 * time.Minute)
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expireAt),
		},
		SessionID: sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func Parse(blacklistStore BlacklistStore, token string) (Claims, error) {
	var claims Claims
	tkn, err := jwt.ParseWithClaims(token, &claims, func(tkn *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		logrus.Errorf("parsing token failed: %v", err)
		if err == jwt.ErrSignatureInvalid {
			return Claims{}, apperror.NewError(apperror.CodeAuthenticationFailed, "parsing token failed")
		}
		// todo this includes token expired error
		return Claims{}, apperror.NewError(apperror.CodeInvalidToken, "invalid token")
	}

	if !tkn.Valid {
		return Claims{}, apperror.NewError(apperror.CodeAuthenticationFailed, "invalid access token")
	}

	isSignedout, err := checkUserIsAlreadySignedOut(token, blacklistStore)
	if err != nil {
		return Claims{}, err
	}
	if isSignedout {
		return Claims{}, apperror.NewError(apperror.CodeAuthenticationFailed, "user has already signed out")
	}

	return claims, nil
}

func checkUserIsAlreadySignedOut(token string, blacklistStore BlacklistStore) (bool, error) {
	isBlacklisted, err := blacklistStore.Exists(token)
	if err != nil {
		// todo return the same error instead of apperror?
		return true, apperror.NewErrorWithCause(apperror.CodeInternalError, "internal error", err)
	}

	if isBlacklisted {
		return true, nil
	}
	return false, nil
}

func BlackList(blacklistStore BlacklistStore, accessToken string) error {
	claims, err := Parse(blacklistStore, accessToken)
	if err != nil {
		return apperror.NewError(apperror.CodeAuthenticationFailed, "")
	}

	return blacklistStore.Add(accessToken, claims.ExpiresAt.Time)
}
