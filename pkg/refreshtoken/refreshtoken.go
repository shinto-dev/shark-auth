package refreshtoken

import (
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

func Create(db *sqlx.DB, userID string, sessionID string) (string, error) {
	userRefreshToken := UserRefreshToken{
		UserID:       userID,
		RefreshToken: uuid.NewV4().String(),
		SessionID:    sessionID,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}
	repo := NewUserRefreshTokenRepository(db)
	err := repo.Create(userRefreshToken)
	if err != nil {
		return "", err
	}

	return userRefreshToken.RefreshToken, nil
}

func Get(db *sqlx.DB, refreshToken string) (UserRefreshToken, error) {
	repo := NewUserRefreshTokenRepository(db)
	return repo.Get(refreshToken)
}

func DeleteBySessionId(db *sqlx.DB, sessionID string) error {
	repo := NewUserRefreshTokenRepository(db)
	return repo.RemoveBySessionID(sessionID)
}

