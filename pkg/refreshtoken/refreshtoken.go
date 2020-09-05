package refreshtoken

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

func Create(tokenStore TokenStore, userID string, sessionID string) (string, error) {
	userRefreshToken := UserRefreshToken{
		UserID:       userID,
		RefreshToken: uuid.NewV4().String(),
		SessionID:    sessionID,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}
	err := tokenStore.Create(userRefreshToken)
	if err != nil {
		return "", err
	}

	return userRefreshToken.RefreshToken, nil
}

func Get(tokenStore TokenStore, refreshToken string) (UserRefreshToken, error) {
	return tokenStore.Get(refreshToken)
}

func DeleteBySessionId(tokenStore TokenStore, sessionID string) error {
	return tokenStore.RemoveBySessionID(sessionID)
}

