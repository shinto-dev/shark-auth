package revoketokens

import (
	"shark-auth/pkg/accesstoken"
	"shark-auth/pkg/refreshtoken"
)

func UsingAccessToken(accessTokenBlackList accesstoken.BlacklistStore, refreshTokenStore refreshtoken.TokenStore, accessToken string) error {
	claims, err := accesstoken.Parse(accessTokenBlackList, accessToken)
	if err != nil {
		return err
	}

	if err = accesstoken.BlackList(accessTokenBlackList, accessToken); err != nil {
		return err
	}

	err = refreshtoken.DeleteBySessionId(refreshTokenStore, claims.SessionID)
	if err != nil {
		return err
	}
	return nil
}
