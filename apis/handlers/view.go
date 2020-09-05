package handlers

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