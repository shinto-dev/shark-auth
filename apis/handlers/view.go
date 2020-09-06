package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/errorcode"
)

type GenericResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   Error       `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code string, message string) GenericResponse {
	return GenericResponse{
		Success: false,
		Data:    nil,
		Error: Error{
			Code:    code,
			Message: message,
		},
	}
}

func NewSuccessResponse(Data interface{}) GenericResponse {
	return GenericResponse{
		Success: true,
		Data:    Data,
	}
}

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

func handleSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // todo: might not be ok always
	if err := json.NewEncoder(w).Encode(NewSuccessResponse(data)); err != nil {
		logrus.Error("writing response json failed")
	}
}

func HandleError(w http.ResponseWriter, err error) {
	errorCode := mapErrorCodeFor(err)
	logrus.WithError(err).
		WithField("stacktrace", fmt.Sprintf("%+v", err)).
		WithField("error_code", errorCode).
		Error(err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mapHttpStatusFor(errorCode))
	if err := json.NewEncoder(w).Encode(NewErrorResponse(errorCode, "")); err != nil {
		logrus.Error("writing response json failed")
	}
}
