package handlers

import (
	"net/http"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"

	"shark-auth/internal/createtokens"
	"shark-auth/internal/revoketokens"
	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/refreshtoken"
	"shark-auth/pkg/user"
)

// This api is for creating new tokens(access token and refresh token) if the user is authenticated.
func HandleTokenCreate(db *sqlx.DB) http.HandlerFunc {
	type GetTokenRequest struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	userRepo := user.NewUserRepository(db)
	refreshTokenRepo := refreshtoken.NewRefreshTokenStore(db)

	return func(w http.ResponseWriter, r *http.Request) {
		var getTokenRequest GetTokenRequest
		if err := readBody(r, &getTokenRequest); err !=nil {
			HandleError(w, apperrors.ErrInvalidJson)
			return
		}

		response, err := createtokens.UsingUserCredentials(userRepo, refreshTokenRepo,
			getTokenRequest.UserName, getTokenRequest.Password)
		if err != nil {
			HandleError(w, err)
			return
		}

		handleSuccess(w, response)
	}
}

func HandleTokenRefresh(db *sqlx.DB) http.HandlerFunc {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	type RefreshTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	refreshTokenStore := refreshtoken.NewRefreshTokenStore(db)

	return func(w http.ResponseWriter, r *http.Request) {
		var refreshTokenRequest RefreshTokenRequest
		if err := readBody(r, &refreshTokenRequest); err != nil {
			HandleError(w, apperrors.ErrInvalidJson)
			return
		}

		jwtToken, err := createtokens.UsingRefreshToken(refreshTokenStore, refreshTokenRequest.RefreshToken)
		if err != nil {
			HandleError(w, err)
			return
		}

		response := RefreshTokenResponse{AccessToken: jwtToken}
		handleSuccess(w, response)
	}
}

func HandleTokenDelete(db *sqlx.DB, redisClient *redis.Client) http.HandlerFunc {
	blacklistStore := accesstoken.NewBlacklistStore(redisClient)
	refreshTokenStore := refreshtoken.NewRefreshTokenStore(db)

	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := extractToken(r)
		if accessToken == "" {
			HandleError(w, apperrors.ErrAccessTokenNotValid)
			return
		}

		err := revoketokens.UsingAccessToken(blacklistStore, refreshTokenStore, accessToken)
		if err != nil {
			HandleError(w, err)
			return
		}

		handleSuccess(w, nil)
	}
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}