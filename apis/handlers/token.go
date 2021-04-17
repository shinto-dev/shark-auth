package handlers

import (
	"net/http"
	"shark-auth/internal"
	"shark-auth/internal/apperror"
	"strings"

	"shark-auth/foundation/logging"
	"shark-auth/foundation/web"
)

type TokenServer struct {
	tokenService internal.TokenService
}

func NewTokenServer(tokenService internal.TokenService) TokenServer {
	server := TokenServer{
		tokenService: tokenService,
	}
	return server
}

// HandleTokenCreate This api is for creating new tokens(access token and refresh token) if the user is authenticated.
func (t *TokenServer) HandleTokenCreate() http.HandlerFunc {
	type GetTokenRequest struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logging.Set(ctx, logging.KeyContext, "create-session")

		var getTokenRequest GetTokenRequest
		if err := readBody(r, &getTokenRequest); err != nil {
			HandleError(ctx, w, apperror.NewError(apperror.CodeInvalidRequest, "invalid json"))
			return
		}

		logging.Set(ctx, "user_name", getTokenRequest.UserName)

		response, err := t.tokenService.CreateToken(getTokenRequest.UserName, getTokenRequest.Password)
		if err != nil {
			HandleError(ctx, w, err)
			return
		}

		web.HandleSuccess(ctx, w, response)
	}
}

func (t *TokenServer) HandleTokenRefresh() http.HandlerFunc {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	type RefreshTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logging.Set(ctx, logging.KeyContext, "refresh-session")

		var refreshTokenRequest RefreshTokenRequest
		if err := readBody(r, &refreshTokenRequest); err != nil {
			HandleError(ctx, w, apperror.NewError(apperror.CodeInvalidRequest, "invalid json"))
			return
		}

		jwtToken, err := t.tokenService.RefreshToken(refreshTokenRequest.RefreshToken)
		if err != nil {
			HandleError(ctx, w, err)
			return
		}

		response := RefreshTokenResponse{AccessToken: jwtToken}
		web.HandleSuccess(r.Context(), w, response)
	}
}

func (t *TokenServer) HandleTokenDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logging.Set(ctx, logging.KeyContext, "delete-session")

		accessToken := extractToken(r)
		if accessToken == "" {
			HandleError(ctx, w, apperror.NewError(apperror.CodeInvalidAccessToken, "access token not valid"))
			return
		}
		err := t.tokenService.Delete(accessToken)
		if err != nil {
			HandleError(ctx, w, err)
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
