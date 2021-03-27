package apperror

import (
	"fmt"
)

type Code string

const (
	CodeAuthenticationFailed Code = "authentication_failed"
	CodeInvalidToken         Code = "invalid_token"
	CodeInternalError        Code = "internal_error"
	CodeUserNotFound         Code = "user_not_found"
	CodeInvalidRequest       Code = "invalid_request"
	CodeInvalidAccessToken   Code = "invalid_access_token"
	CodePasswordMismatch     Code = "password_mismatch"
	CodeInvalidRefreshToken  Code = "invalid_refresh_token"
)

type Error struct {
	Code  Code
	Msg   string
	Cause error
}

func NewErrorWithCause(code Code, msg string, cause error) *Error {
	return &Error{Code: code, Msg: msg, Cause: cause}
}

func NewError(code Code, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

func (a Error) Error() string {
	return fmt.Sprintf("errorcode: %s, message: %s:%s", a.Code, a.Msg, a.Cause)
}
