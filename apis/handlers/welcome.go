package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"shark-auth/autherrors"
	"shark-auth/pkg/token"
)

func Welcome(c *gin.Context) {
	accessToken := extractToken(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, "token not valid")
	}

	claims, err := token.ParseAccessToken(accessToken)
	if err != nil {
		if err == autherrors.ErrAuthenticationFailed {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}
	logrus.Infof("request received from user: %s", claims.Username)

	c.Writer.Write([]byte("Hello world"))
}
