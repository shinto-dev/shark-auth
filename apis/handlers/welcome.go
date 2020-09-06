package handlers

import (
	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"

	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/apperrors"
)

// a sample api which can be used for testing the authentication
func HandleWelcome(redisClient *redis.Client) http.HandlerFunc {
	blacklistStore := accesstoken.NewBlacklistStore(redisClient)

	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := extractToken(r)
		if accessToken == "" {
			HandleError(w, apperrors.ErrAccessTokenNotValid)
			return
		}

		claims, err := accesstoken.Parse(blacklistStore, accessToken)
		if err != nil {
			HandleError(w, err)
			return
		}
		logrus.Infof("request received from user: %s", claims.UserID)

		handleSuccess(w, "Hello world")
	}
}
