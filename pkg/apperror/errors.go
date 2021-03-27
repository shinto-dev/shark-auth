package apperror

import (
	"errors"
)

var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrInvalidToken = errors.New("invalid token")
var ErrInternal = errors.New("internal error")
var ErrPasswordMismatch = errors.New("password mismatch")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidJson = errors.New("invalid json input")
var ErrRefreshTokenNotValid = errors.New("refresh token not valid")
var ErrAccessTokenNotValid = errors.New("access token not valid")
var ErrUserNameNotAvailable = errors.New("user name already taken")
