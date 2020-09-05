package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"shark-auth/pkg/apperrors"
	"shark-auth/pkg/errorcode"
)

type GenericResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   Error       `json:"error"`
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

func handleError(c *gin.Context, err error) {
	errorCode := mapErrorCodeFor(err)
	logrus.WithError(err).
		WithField("error_code", errorCode).
		Error("error while creating tokens")
	c.JSON(mapHttpStatusFor(errorCode), NewErrorResponse(errorCode, ""))
}

func mapErrorCodeFor(err error) string {
	switch err {
	case apperrors.ErrPasswordMismatch:
		return errorcode.ERROR_AUTHENTICATION_FAILED
	case apperrors.ErrUserNotFound:
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
