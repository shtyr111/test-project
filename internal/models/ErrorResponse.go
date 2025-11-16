package models

type ErrorDetail struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorDetail ErrorDetail `json:"error"`
}

func NewErrorResponse(Status int, Message string) *ErrorResponse {
	errorDetail := ErrorDetail{Status: Status, Message: Message}
	return &ErrorResponse{ErrorDetail: errorDetail}
}
