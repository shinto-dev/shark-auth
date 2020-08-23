package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"shark-auth/autherrors"
	"shark-auth/pkg/accesstoken"
)

func DeleteToken(c *gin.Context) {
	accessToken := extractToken(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, "token not valid")
	}

	err := accesstoken.DeleteAccessToken(accessToken)
	if err != nil {
		if err == autherrors.ErrAuthenticationFailed {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	// todo add session id and remove the refresh token

	c.Status(http.StatusOK)
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
