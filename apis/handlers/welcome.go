package handlers

import (
	"net/http"
	"shark-auth/internal/accesstoken"
	"shark-auth/internal/apperror"

	"github.com/sirupsen/logrus"

	"shark-auth/foundation/web"
)

//HandleWelcome is a sample api which can be used for testing the authentication
func HandleWelcome(blacklistStore accesstoken.BlacklistStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		accessToken := extractToken(r)
		if accessToken == "" {
			HandleError(ctx, w, apperror.NewError(apperror.CodeInvalidAccessToken, "access token not valid"))
			return
		}

		claims, err := accesstoken.Parse(blacklistStore, accessToken)
		if err != nil {
			HandleError(ctx, w, err)
			return
		}
		logrus.Infof("request received from user: %s", claims.UserID)

		web.HandleSuccess(r.Context(), w, "Hello world")
	}
}
