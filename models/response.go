package models

// Response represents a generic response message
type Response struct {
	Message string `json:"message"`
}

func NewResponse(message string) Response {
	return Response{
		Message: message,
	}
}

func NewErrorResponse(message string, err string) Response {
	if message == "" || err == "" {
		return Response{
			Message: message + err,
		}
	}

	return Response{
		Message: message + ": " + err,
	}
}