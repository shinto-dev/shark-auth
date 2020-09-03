package accesstoken

import (
	"fmt"
	"time"

	"shark-auth/foundation/redis_client"
)

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

func blacklistAccessTokenKey(refreshToken string) string {
	return fmt.Sprintf("blacklist-access-token:%s", refreshToken)
}
