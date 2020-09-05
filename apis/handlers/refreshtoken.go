package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"shark-auth/internal/createtokens"
	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/refreshtoken"
)

func HandleTokenRefresh(db *sqlx.DB) func(c *gin.Context) {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	type RefreshTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	refreshTokenStore := refreshtoken.NewRefreshTokenStore(db)

	return func(c *gin.Context) {
		var refreshTokenRequest RefreshTokenRequest
		if err := c.BindJSON(&refreshTokenRequest); err != nil {
			handleError(c, apperrors.ErrInvalidJson)
			return
		}

		jwtToken, err := createtokens.UsingRefreshToken(refreshTokenStore, refreshTokenRequest.RefreshToken)
		if err != nil {
			handleError(c, err)
			return
		}

		response := RefreshTokenResponse{AccessToken: jwtToken}
		c.JSON(http.StatusOK, NewSuccessResponse(response))
	}
}
