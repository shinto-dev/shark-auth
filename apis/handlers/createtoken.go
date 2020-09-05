package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"shark-auth/internal/createtokens"
	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/refreshtoken"
	"shark-auth/pkg/user"
)

func HandleTokenCreate(db *sqlx.DB) func(c *gin.Context) {
	type GetTokenRequest struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	userRepo := user.NewUserRepository(db)
	refreshTokenRepo := refreshtoken.NewRefreshTokenStore(db)

	return func(c *gin.Context) {
		var getTokenRequest GetTokenRequest
		if err := c.ShouldBindJSON(&getTokenRequest); err != nil {
			handleError(c, apperrors.ErrInvalidJson)
			return
		}

		response, err := createtokens.UsingUserCredentials(userRepo, refreshTokenRepo,
			getTokenRequest.UserName, getTokenRequest.Password)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, NewSuccessResponse(response))
	}
}
