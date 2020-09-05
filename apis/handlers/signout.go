package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"

	"shark-auth/autherrors"
	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/refreshtoken"
)

func DeleteToken(db *sqlx.DB, redisClient *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		accessToken := extractToken(c)
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, "token not valid")
		}

		claims, err := accesstoken.Parse(accessToken, redisClient)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
			// todo: handle generic errors if any
		}

		err = accesstoken.Delete(accessToken, redisClient)
		if err != nil {
			if err == autherrors.ErrAuthenticationFailed {
				c.Status(http.StatusUnauthorized)
				return
			}
			c.Status(http.StatusBadRequest)
			return
		}

		err = refreshtoken.DeleteBySessionId(db, claims.SessionID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
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
