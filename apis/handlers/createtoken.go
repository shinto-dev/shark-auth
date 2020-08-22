package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"shark-auth/pkg/token"
	"shark-auth/pkg/user"
)

type GetTokenRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type GetTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GetToken(userRepo user.UserRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		var getTokenRequest GetTokenRequest
		if err := c.ShouldBindJSON(&getTokenRequest); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid json provided")
			return
		}

		validUser, err := userRepo.IsValid(getTokenRequest.UserName, getTokenRequest.Password)
		// todo handle error or panic from calling function
		if !validUser {
			logrus.WithField("user_name", getTokenRequest.UserName).
				Error("password does not match")
			c.Status(http.StatusUnauthorized)
			return
		}

		tkn, err := token.CreateAccessToken(getTokenRequest.UserName)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		refreshTkn, err := token.CreateRefreshToken(getTokenRequest.UserName)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		response := GetTokenResponse{
			AccessToken:  tkn,
			RefreshToken: refreshTkn,
		}
		c.JSON(http.StatusOK, response)
	}
}
