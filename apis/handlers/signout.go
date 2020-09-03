package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"shark-auth/autherrors"
	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/refreshtoken"
)

func DeleteToken(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		accessToken := extractToken(c)
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, "token not valid")
		}

		claims, err := accesstoken.ParseAccessToken(accessToken)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
			// todo: handle generic errors if any
		}

		err = accesstoken.DeleteAccessToken(accessToken)
		if err != nil {
			if err == autherrors.ErrAuthenticationFailed {
				c.Status(http.StatusUnauthorized)
				return
			}
			c.Status(http.StatusBadRequest)
			return
		}

		err = refreshtoken.DeleteRefreshTokenBySessionId(db, claims.SessionID)
		if err!=nil {
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
