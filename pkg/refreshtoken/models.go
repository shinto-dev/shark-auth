package refreshtoken

import (
	"time"
)

type UserRefreshToken struct {
	UserID       string    `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	SessionID    string    `db:"session_id"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}
