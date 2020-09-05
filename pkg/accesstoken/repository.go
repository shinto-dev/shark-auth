package accesstoken

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type AccessTokenBlacklistStore interface {
	Add(accessToken string, expiresAt time.Time) error
	Exists(accessToken string) (bool, error)
}

func BlacklistAccessToken(accessToken string, expiresAt time.Time, redisClient *redis.Client) error {
	return redisClient.Set(blacklistAccessTokenKey(accessToken), true, expiresAt.Sub(time.Now())).Err()
}

func IsAccessTokenBlacklisted(accessToken string, redisClient *redis.Client) (bool, error) {
	cmd := redisClient.Exists(blacklistAccessTokenKey(accessToken))
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	// todo recheck this
	return cmd.Val() == 1, nil
}

func blacklistAccessTokenKey(refreshToken string) string {
	return fmt.Sprintf("blacklist-access-token:%s", refreshToken)
}
