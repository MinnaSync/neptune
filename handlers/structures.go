package handlers

type APIResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type APIError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
