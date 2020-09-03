package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"shark-auth/autherrors"
	"shark-auth/pkg/accesstoken"
)

func Welcome(c *gin.Context) {
	accessToken := extractToken(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, "token not valid")
	}

	claims, err := accesstoken.Parse(accessToken)
	if err != nil {
		if err == autherrors.ErrAuthenticationFailed {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}
	logrus.Infof("request received from user: %s", claims.UserID)

	c.Writer.Write([]byte("Hello world"))
}
