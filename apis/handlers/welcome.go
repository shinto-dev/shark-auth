package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"shark-auth/foundation/web"
	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/apperror"
)

//HandleWelcome is a sample api which can be used for testing the authentication
func HandleWelcome(blacklistStore accesstoken.BlacklistStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := extractToken(r)
		if accessToken == "" {
			HandleError(w, apperror.NewError(apperror.CodeInvalidAccessToken, "access token not valid"))
			return
		}

		claims, err := accesstoken.Parse(blacklistStore, accessToken)
		if err != nil {
			HandleError(w, err)
			return
		}
		logrus.Infof("request received from user: %s", claims.UserID)

		web.HandleSuccess(r.Context(), w, "Hello world")
	}
}
