package apis

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"shark-auth/autherrors"
	"shark-auth/token"
	"shark-auth/user"
)

const TOKEN = "token"

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

		if !userRepo.IsValid(getTokenRequest.UserName, getTokenRequest.Password) {
			logrus.Error("password does not match")
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

func DeleteToken(c *gin.Context) {
	accessToken := extractToken(c)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, "token not valid")
	}

	err := token.DeleteAccessToken(accessToken)
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

func welcome(c *gin.Context) {
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
