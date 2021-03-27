package handlers

import (
	"net/http"

	"shark-auth/foundation/web"
	"shark-auth/internal/signupuser"
	"shark-auth/pkg/apperror"
	"shark-auth/pkg/user"
)

//HandleUserSignup is a very basic signup api
func HandleUserSignup(userRepo user.Repository) http.HandlerFunc {
	type SignupRequest struct {
		UserName string
		Password string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var signupRequest SignupRequest
		if err := readBody(r, &signupRequest); err != nil {
			HandleError(w, apperror.NewErrorWithCause(apperror.CodeInvalidRequest, "invalid json", err))
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

		web.HandleSuccess(r.Context(), w, nil)
	}
}
