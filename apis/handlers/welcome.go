package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"

	"shark-auth/pkg/accesstoken"
)

func HandleWelcome(redisClient *redis.Client) func(c *gin.Context) {
	blacklistStore := accesstoken.NewBlacklistStore(redisClient)

	return func(c *gin.Context) {
		accessToken := extractToken(c)
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, "token not valid")
		}

		claims, err := accesstoken.Parse(blacklistStore, accessToken)
		if err != nil {
			handleError(c, err)
			return
		}
		logrus.Infof("request received from user: %s", claims.UserID)

		c.Writer.Write([]byte("Hello world"))
	}
}
