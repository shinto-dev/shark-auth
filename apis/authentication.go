package apis

import (
	"net/http"
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

		tkn, err := token.CreateJwtToken(getTokenRequest.UserName)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		response := GetTokenResponse{AccessToken: tkn}
		c.JSON(http.StatusOK, response)
	}
}

func Refresh(c *gin.Context) {
	tknStr, err := c.Cookie(TOKEN)
	if err != nil {
		if err == http.ErrNoCookie {
			logrus.Error("cookie not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	claims, err := token.ParseJwtToken(tknStr)
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
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.Status(http.StatusBadRequest)
		return
	}

	jwtToken, err := token.CreateJwtToken(claims.Username)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	response := GetTokenResponse{AccessToken: jwtToken}
	c.JSON(http.StatusOK, response)
}


func welcome(c *gin.Context) {
	tokenString, err := c.Cookie(TOKEN)
	if err != nil {
		if err == http.ErrNoCookie {
			logrus.Error("cookie not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	claims, err := token.ParseJwtToken(tokenString)
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