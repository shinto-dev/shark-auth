package refreshtoken

import (
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

// todo move this to postgres repo
func CreateRefreshToken(db *sqlx.DB, userID string) (string, error) {
	userRefreshToken := UserRefreshToken{
		UserID:       userID,
		RefreshToken: uuid.NewV4().String(),
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

func IsRefreshTokenValid(db *sqlx.DB, refreshToken string, userID string) (bool, error) {
	repo := NewUserRefreshTokenRepository(db)
	userRefreshToken, err := repo.Get(refreshToken)
	if err != nil {
		return false, err
	}

	//todo return userRefreshToken.UserID == userID, nil
	return userRefreshToken != (UserRefreshToken{}), nil
}

func DeleteRefreshToken(db *sqlx.DB, refreshToken string) error {
	// todo include this: either accept session id or refresh token
	repo := NewUserRefreshTokenRepository(db)
	return repo.Remove(refreshToken)
}