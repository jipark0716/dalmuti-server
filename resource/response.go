package resource

type InvalidRequestResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewInvalidRequestResponse(message string) InvalidRequestResponse {
	return InvalidRequestResponse{
		Code:    422,
		Message: message,
	}
}
