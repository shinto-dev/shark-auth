package createtokens

import (
	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/refreshtoken"
)

func UsingRefreshToken(store refreshtoken.TokenStore, refreshToken string) (string, error) {
	userRefreshToken, err := refreshtoken.Get(store, refreshToken)
	if err != nil {
		return "", err
	}

	if userRefreshToken == (refreshtoken.UserRefreshToken{}) {
		return "", apperrors.ErrRefreshTokenNotValid
	}

	jwtToken, err := accesstoken.Create(userRefreshToken.UserID, userRefreshToken.SessionID)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
