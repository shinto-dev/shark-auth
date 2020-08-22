package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"shark-auth/autherrors"
	"shark-auth/pkg/refreshtoken"
	"shark-auth/pkg/token"
)

func RefreshToken(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		accessToken := extractToken(c)
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, "token not valid")
			return
		}
		refreshToken := c.GetHeader("refresh-token")

		claims, err := accesstoken.ParseAccessToken(accessToken)
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
			logrus.Error("token refresh is not allowed")
			c.Status(http.StatusBadRequest)
			return
		}

		var tokenValid bool
		if tokenValid, err = refreshtoken.IsRefreshTokenValid(db, refreshToken, claims.Username); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if !tokenValid {
			c.JSON(http.StatusUnauthorized, "refresh token not valid")
			return
		}

		jwtToken, err := accesstoken.CreateAccessToken(claims.Username)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		response := GetTokenResponse{AccessToken: jwtToken}
		c.JSON(http.StatusOK, response)
	}
}
