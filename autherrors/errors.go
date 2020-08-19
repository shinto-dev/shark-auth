package autherrors

import (
	"errors"
)

const AuthenticationFailed = "002"
const InvalidToken = "003"

var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrInvalidToken = errors.New("invalid token")
var ErrInternal = errors.New("internal error")