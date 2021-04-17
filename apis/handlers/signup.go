package handlers

import (
	"net/http"
	"shark-auth/internal"
	"shark-auth/internal/apperror"
	"shark-auth/internal/user"

	"shark-auth/foundation/web"
)

//HandleUserSignup is a very basic signup api
func HandleUserSignup(userRepo user.Repository) http.HandlerFunc {
	type SignupRequest struct {
		UserName string
		Password string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var signupRequest SignupRequest
		if err := readBody(r, &signupRequest); err != nil {
			HandleError(ctx, w, apperror.NewErrorWithCause(apperror.CodeInvalidRequest, "invalid json", err))
			return
		}

		// todo validations
		userDetails := internal.User{
			UserName: signupRequest.UserName,
			Password: signupRequest.Password,
		}

		if err := internal.CreateUser(userRepo, userDetails); err != nil {
			HandleError(ctx, w, err)
			return
		}

		web.HandleSuccess(r.Context(), w, nil)
	}
}
