package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"shark-auth/foundation/logging"
	"shark-auth/internal/apperror"

	"shark-auth/foundation/web"
)

func mapErrorCodeFor(err error) apperror.Code {
	appErr, ok := err.(apperror.Error)
	if !ok {
		return apperror.CodeInternalError
	}

	return appErr.Code
}

func mapHttpStatusFor(errorCode apperror.Code) int {
	switch errorCode {
	case apperror.CodeAuthenticationFailed:
		return http.StatusUnauthorized
	case apperror.CodeInvalidRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func readBody(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error) {
	errorCode := mapErrorCodeFor(err)
	logging.FromContext(ctx).WithError(err).
		WithField("stacktrace", fmt.Sprintf("%+v", err)).
		WithField("error_code", errorCode).
		Error(err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mapHttpStatusFor(errorCode))
	if err := json.NewEncoder(w).Encode(web.NewErrorResponse(errorCode, "")); err != nil {
		logging.FromContext(ctx).Error("writing response json failed")
	}
}
