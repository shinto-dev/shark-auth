package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"shark-auth/internal/signupuser"
	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/user"
)

// this is a very basic signup api
func HandleUserSignup(db *sqlx.DB) func(c *gin.Context) {
	userRepo := user.NewUserRepository(db)

	type SignupRequest struct {
		UserName string
		Password string
	}

	return func(c *gin.Context) {
		var signupRequest SignupRequest
		if err := c.ShouldBindJSON(&signupRequest); err != nil {
			handleError(c, apperrors.ErrInvalidJson)
			return
		}

		// todo validations
		userDetails := signupuser.User{
			UserName: signupRequest.UserName,
			Password: signupRequest.Password,
		}

		if err := signupuser.CreateUser(userRepo, userDetails); err != nil {
			handleError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
