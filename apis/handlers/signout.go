package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"

	"shark-auth/internal/revoketokens"
	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/refreshtoken"
)

func HandleTokenDelete(db *sqlx.DB, redisClient *redis.Client) func(c *gin.Context) {
	blacklistStore := accesstoken.NewBlacklistStore(redisClient)
	refreshTokenStore := refreshtoken.NewRefreshTokenStore(db)

	return func(c *gin.Context) {
		accessToken := extractToken(c)
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, "token not valid")
			return
		}

		err := revoketokens.UsingAccessToken(blacklistStore, refreshTokenStore, accessToken)
		if err != nil {
			handleError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
