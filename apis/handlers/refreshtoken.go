package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/refreshtoken"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshToken(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var refreshTokenRequest RefreshTokenRequest
		if err := c.BindJSON(&refreshTokenRequest); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		userRefreshToken, err := refreshtoken.Get(db, refreshTokenRequest.RefreshToken)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if userRefreshToken == (refreshtoken.UserRefreshToken{}) {
			c.JSON(http.StatusUnauthorized, "refresh token not valid")
			return
		}

		jwtToken, err := accesstoken.Create(userRefreshToken.UserID, userRefreshToken.SessionID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		response := GetTokenResponse{AccessToken: jwtToken}
		c.JSON(http.StatusOK, response)
	}
}
