package handlers

import (
	"net/http"
	"strings"

	"shark-auth/foundation/logging"
	"shark-auth/foundation/web"
	"shark-auth/internal/createtokens"
	"shark-auth/internal/revoketokens"
	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/apperror"
	"shark-auth/pkg/refreshtoken"
	"shark-auth/pkg/user"
)

type TokenServer struct {
	userRepo                  user.Repository
	refreshTokenStore         refreshtoken.TokenStore
	accessTokenBlacklistStore accesstoken.BlacklistStore
}

func NewTokenServer(UserRepo user.Repository, RefreshTokenStore refreshtoken.TokenStore,
	AccessTokenBlacklistStore accesstoken.BlacklistStore) TokenServer {
	server := TokenServer{
		userRepo:                  UserRepo,
		refreshTokenStore:         RefreshTokenStore,
		accessTokenBlacklistStore: AccessTokenBlacklistStore,
	}
	return server
}

// This api is for creating new tokens(access token and refresh token) if the user is authenticated.
func (t TokenServer) HandleTokenCreate() http.HandlerFunc {
	type GetTokenRequest struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var getTokenRequest GetTokenRequest
		if err := readBody(r, &getTokenRequest); err != nil {
			HandleError(w, apperror.ErrInvalidJson)
			return
		}

		ctx := r.Context()
		logging.Set(ctx, "user_name", getTokenRequest.UserName)

		response, err := createtokens.UsingUserCredentials(t.userRepo, t.refreshTokenStore,
			getTokenRequest.UserName, getTokenRequest.Password)
		if err != nil {
			HandleError(w, err)
			return
		}

		web.HandleSuccess(ctx, w, response)
	}
}

func (t TokenServer) HandleTokenRefresh() http.HandlerFunc {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	type RefreshTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var refreshTokenRequest RefreshTokenRequest
		if err := readBody(r, &refreshTokenRequest); err != nil {
			HandleError(w, apperror.ErrInvalidJson)
			return
		}

		jwtToken, err := createtokens.UsingRefreshToken(t.refreshTokenStore, refreshTokenRequest.RefreshToken)
		if err != nil {
			HandleError(w, err)
			return
		}

		response := RefreshTokenResponse{AccessToken: jwtToken}
		web.HandleSuccess(r.Context(), w, response)
	}
}

func (t TokenServer) HandleTokenDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := extractToken(r)
		if accessToken == "" {
			HandleError(w, apperror.ErrAccessTokenNotValid)
			return
		}

		err := revoketokens.UsingAccessToken(t.accessTokenBlacklistStore, t.refreshTokenStore, accessToken)
		if err != nil {
			HandleError(w, err)
			return
		}

		web.HandleSuccess(r.Context(), w, nil)
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
