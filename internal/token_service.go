package internal

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"shark-auth/internal/accesstoken"
	"shark-auth/internal/apperror"
	"shark-auth/internal/refreshtoken"
	"shark-auth/internal/user"
)

type TokenService struct {
	UserRepo          user.Repository
	RefreshTokenStore refreshtoken.TokenStore
	BlacklistStore    accesstoken.BlacklistStore
}

type CreateTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (t *TokenService) CreateToken(userName, password string) (CreateTokenResponse, error) {
	currentUser, err := user.GetByUserNameAndPassword(t.UserRepo, userName, password)
	if err != nil {
		return CreateTokenResponse{}, err
	}

	return t.createTokenFor(currentUser)
}

func (t *TokenService) RefreshToken(refreshToken string) (string, error) {
	userRefreshToken, err := refreshtoken.Get(t.RefreshTokenStore, refreshToken)
	if err != nil {
		return "", err
	}

	if userRefreshToken == (refreshtoken.UserRefreshToken{}) {
		return "", apperror.NewError(apperror.CodeInvalidRefreshToken, "refresh token is invalid")
	}

	jwtToken, err := accesstoken.Create(userRefreshToken.UserID, userRefreshToken.SessionID)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (t *TokenService) Delete(accessToken string) error {
	claims, err := accesstoken.Parse(t.BlacklistStore, accessToken)
	if err != nil {
		return err
	}

	if err = accesstoken.BlackList(t.BlacklistStore, accessToken); err != nil {
		return err
	}

	err = refreshtoken.DeleteBySessionId(t.RefreshTokenStore, claims.SessionID)
	if err != nil {
		return err
	}
	return nil
}

// todo more session details, device info(or browser info)?
func (t *TokenService) createTokenFor(currentUser user.User) (CreateTokenResponse, error) {
	sessionID := uuid.NewV4().String()

	tkn, err := accesstoken.Create(currentUser.UserId, sessionID)
	if err != nil {
		return CreateTokenResponse{}, errors.Wrap(err, "error while creating access token")
	}

	refreshTkn, err := refreshtoken.Create(t.RefreshTokenStore, currentUser.UserId, sessionID)
	if err != nil {
		return CreateTokenResponse{}, errors.Wrap(err, "error while creating refresh token")
	}

	return CreateTokenResponse{
		AccessToken:  tkn,
		RefreshToken: refreshTkn,
	}, nil
}
