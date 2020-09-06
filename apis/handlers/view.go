package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"shark-auth/foundation/web"
	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/errorcode"
)

func mapErrorCodeFor(err error) string {
	switch err {
	case apperrors.ErrPasswordMismatch:
		return errorcode.ERROR_AUTHENTICATION_FAILED
	case apperrors.ErrUserNotFound:
		return errorcode.ERROR_AUTHENTICATION_FAILED
	case apperrors.ErrAccessTokenNotValid:
		return errorcode.ERROR_AUTHENTICATION_FAILED
	case apperrors.ErrInvalidToken:
		return errorcode.ERROR_BAD_REQUEST
	case apperrors.ErrInvalidJson:
		return errorcode.ERROR_BAD_REQUEST
	case apperrors.ErrUserNameNotAvailable:
		return errorcode.USERNAME_NOT_AVAILABLE
	default:
		return errorcode.ERROR_INTERNAL
	}
}

func mapHttpStatusFor(errorCode string) int {
	switch errorCode {
	case errorcode.ERROR_AUTHENTICATION_FAILED:
		return http.StatusUnauthorized
	case errorcode.ERROR_BAD_REQUEST:
		return http.StatusBadRequest
	case errorcode.USERNAME_NOT_AVAILABLE:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func readBody(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func HandleError(w http.ResponseWriter, err error) {
	errorCode := mapErrorCodeFor(err)
	logrus.WithError(err).
		WithField("stacktrace", fmt.Sprintf("%+v", err)).
		WithField("error_code", errorCode).
		Error(err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mapHttpStatusFor(errorCode))
	if err := json.NewEncoder(w).Encode(web.NewErrorResponse(errorCode, "")); err != nil {
		logrus.Error("writing response json failed")
	}
}
