package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/errorcode"
	"shark-auth/pkg/refreshtoken"
	"shark-auth/pkg/user"
)

type GetTokenRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type GetTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func CreateToken(userRepo user.UserRepository, db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var getTokenRequest GetTokenRequest
		if err := c.ShouldBindJSON(&getTokenRequest); err != nil {
			response := NewErrorResponse(errorcode.ERROR_BAD_REQUEST, "invalid json provided")
			c.JSON(http.StatusBadRequest, response)
			return
		}

		currentUser, err := userRepo.Get(getTokenRequest.UserName)
		if err != nil {
			// todo panic from calling function
			logrus.WithError(err).Error("error while retrieving user")
			c.Status(http.StatusInternalServerError)
			return
		}

		if currentUser == (user.User{}) || !passwordMatch(currentUser.Password, getTokenRequest.Password) {
			logrus.WithField("user_name", getTokenRequest.UserName).
				Error("password does not match")
			response := NewErrorResponse(errorcode.ERROR_PASSWORD_MISMATCH, "password does not match")
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		// todo more session details, device info(or browser info)?
		sessionID := uuid.NewV4().String()

		tkn, err := accesstoken.Create(currentUser.UserId, sessionID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		refreshTkn, err := refreshtoken.Create(db, currentUser.UserId, sessionID)
		if err != nil {
			logrus.WithError(err).Error("refresh token creation failed")
			c.Status(http.StatusInternalServerError)
			return
		}

		response := GetTokenResponse{
			AccessToken:  tkn,
			RefreshToken: refreshTkn,
		}
		c.JSON(http.StatusOK, NewSuccessResponse(response))
	}
}

func passwordMatch(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
