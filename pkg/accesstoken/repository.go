package accesstoken

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type BlacklistStore interface {
	Add(accessToken string, expiresAt time.Time) error
	Exists(accessToken string) (bool, error)
}

func NewBlacklistStore(redisClient *redis.Client) BlacklistStore {
	return &AccessTokenBlacklist{
		redisClient: redisClient,
	}
}

type AccessTokenBlacklist struct {
	redisClient *redis.Client
}

func (a *AccessTokenBlacklist) Add(accessToken string, expiresAt time.Time) error {
	return a.redisClient.Set(blacklistAccessTokenKey(accessToken), true, expiresAt.Sub(time.Now())).Err()
}

func (a *AccessTokenBlacklist) Exists(accessToken string) (bool, error) {
	cmd := a.redisClient.Exists(blacklistAccessTokenKey(accessToken))
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	// todo recheck this
	return cmd.Val() == 1, nil
}

func blacklistAccessTokenKey(refreshToken string) string {
	return fmt.Sprintf("blacklist-access-token:%s", refreshToken)
}
