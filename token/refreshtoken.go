package token

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// todo move this to postgres repo
func CreateRefreshToken(userName string) (string, error) {
	expireAt := time.Now().Add(7 * 24 * time.Hour)
	refreshToken := uuid.NewV4().String()
	if err := SaveRefreshToken(userName, refreshToken, expireAt); err != nil {
		return "", err
	}

	return refreshToken, nil
}

func IsRefreshTokenValid(refreshToken string, userName string) (bool, error) {
	refreshTokenUserName, err := GetRefreshToken(refreshToken)
	if err != nil {
		return false, err
	}

	return refreshTokenUserName == userName, nil
}
