package web

import (
	"context"
	"encoding/json"
	"net/http"
	"shark-auth/foundation/logging"
	"shark-auth/internal/apperror"
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

//NewSuccessResponse creates new success response
func NewSuccessResponse(Data interface{}) GenericResponse {
	return GenericResponse{
		Success: true,
		Data:    Data,
	}
}

//NewErrorResponse creates error response
func NewErrorResponse(code apperror.Code, message string) GenericResponse {
	return GenericResponse{
		Success: false,
		Data:    nil,
		Error: Error{
			Code:    string(code),
			Message: message,
		},
	}
}

func HandleSuccess(ctx context.Context, w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // todo: might not be ok always

	if err := json.NewEncoder(w).Encode(NewSuccessResponse(data)); err != nil {
		logging.FromContext(ctx).Error("writing response json failed")
	}
}
