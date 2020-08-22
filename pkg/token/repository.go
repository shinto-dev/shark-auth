package token

import (
	"fmt"
	"time"

	"shark-auth/foundation/redis_client"
)

func SaveRefreshToken(userID string, refreshToken string, expiresAt time.Time) error {
	return redis_client.Client.Set(refreshTokenKey(refreshToken), userID, expiresAt.Sub(time.Now())).Err()
}

func GetRefreshToken(refreshToken string) (string, error) {
	cmd := redis_client.Client.Get(refreshTokenKey(refreshToken))
	if cmd.Err() != nil {
		return "", cmd.Err()
	}

	return cmd.Val(), nil
}

func RemoveRefreshToken(refreshToken string) error {
	return redis_client.Client.Unlink(refreshTokenKey(refreshToken)).Err()
}

func BlacklistAccessToken(accessToken string, expiresAt time.Time) error {
	return redis_client.Client.Set(blacklistAccessTokenKey(accessToken), true, expiresAt.Sub(time.Now())).Err()
}

func IsAccessTokenBlacklisted(accessToken string) (bool, error) {
	cmd := redis_client.Client.Exists(blacklistAccessTokenKey(accessToken))
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	// todo recheck this
	return cmd.Val() == 1, nil
}

func refreshTokenKey(refreshToken string) string {
	return fmt.Sprintf("refresh-token:%s", refreshToken)
}

func blacklistAccessTokenKey(refreshToken string) string {
	return fmt.Sprintf("blacklist-access-token:%s", refreshToken)
}