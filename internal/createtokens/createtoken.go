package createtokens

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/refreshtoken"
	"shark-auth/pkg/user"
)

type CreateTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func UsingUserCredentials(userRepo user.Repository, refreshTokenStore refreshtoken.TokenStore,
	userName, password string) (CreateTokenResponse, error) {
	currentUser, err := user.GetByUserNameAndPassword(userRepo, userName, password)
	if err != nil {
		return CreateTokenResponse{}, err
	}

	return createTokenFor(refreshTokenStore, currentUser)
}

// todo more session details, device info(or browser info)?
func createTokenFor(refreshTokenStore refreshtoken.TokenStore,
	currentUser user.User) (CreateTokenResponse, error) {
	sessionID := uuid.NewV4().String()

	tkn, err := accesstoken.Create(currentUser.UserId, sessionID)
	if err != nil {
		return CreateTokenResponse{}, errors.Wrap(err, "error while creating access token")
	}

	refreshTkn, err := refreshtoken.Create(refreshTokenStore, currentUser.UserId, sessionID)
	if err != nil {
		return CreateTokenResponse{}, errors.Wrap(err, "error while creating refresh token")
	}

	return CreateTokenResponse{
		AccessToken:  tkn,
		RefreshToken: refreshTkn,
	}, nil
}
