package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"shark-auth/autherrors"
	"shark-auth/pkg/token"
)

func Refresh(c *gin.Context) {
	accessToken := extractToken(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, "token not valid")
	}
	refreshToken := c.GetHeader("refresh-token")

	claims, err := token.ParseAccessToken(accessToken)
	if err != nil {
		if err == autherrors.ErrAuthenticationFailed {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new tkn will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if claims.ExpiresAt.Sub(time.Now()) > 30*time.Second {
		c.Status(http.StatusBadRequest)
		return
	}

	var tokenValid bool
	if tokenValid, err = token.IsRefreshTokenValid(refreshToken, claims.Username); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if !tokenValid {
		c.JSON(http.StatusUnauthorized, "refresh token not valid")
		return
	}

	jwtToken, err := token.CreateAccessToken(claims.Username)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	response := GetTokenResponse{AccessToken: jwtToken}
	c.JSON(http.StatusOK, response)
}
