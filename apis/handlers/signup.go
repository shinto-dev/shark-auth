package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"shark-auth/pkg/user"
)

type SignupRequest struct {
	UserName string
	Password string
}

// this is a very basic signup api
func Signup(userRepo user.UserRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		// todo validations
		var signupRequest SignupRequest
		if err := c.ShouldBindJSON(&signupRequest); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid json provided")
			return
		}

		currentUser, err := userRepo.Get(signupRequest.UserName)
		if err != nil {
			// todo panic from calling function
			logrus.WithError(err).Error("error while retrieving user")
			c.Status(http.StatusInternalServerError)
			return
		}

		if currentUser != (user.User{}) {
			logrus.WithField("user_name", signupRequest.UserName).
				Error("user name already taken")
			c.Status(http.StatusBadRequest)
			return
		}

		err = userRepo.Create(user.User{
			UserId:    uuid.NewV4().String(),
			UserName:  signupRequest.UserName,
			Password:  hashPassword(signupRequest.Password),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			logrus.WithError(err).Error("error inserting user")
			c.Status(http.StatusInternalServerError)
		}

		c.Status(http.StatusOK)
	}
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	// todo handle this error
	return string(hashedPassword)
}
