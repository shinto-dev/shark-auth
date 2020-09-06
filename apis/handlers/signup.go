package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"shark-auth/internal/signupuser"
	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/user"
)

// this is a very basic signup api
func HandleUserSignup(db *sqlx.DB) http.HandlerFunc {
	userRepo := user.NewUserRepository(db)

	type SignupRequest struct {
		UserName string
		Password string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var signupRequest SignupRequest
		if err := readBody(r, signupRequest); err != nil {
			HandleError(w, apperrors.ErrInvalidJson)
			return
		}

		// todo validations
		userDetails := signupuser.User{
			UserName: signupRequest.UserName,
			Password: signupRequest.Password,
		}

		if err := signupuser.CreateUser(userRepo, userDetails); err != nil {
			HandleError(w, err)
			return
		}

		handleSuccess(w, nil)
	}
}
