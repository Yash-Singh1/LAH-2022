package utils

type APIError struct {
	ErrorMessage string `json:"error"`
}

type ErrorResponse struct {
	Code int
	Body APIError
}
