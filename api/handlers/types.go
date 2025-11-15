package handlers

type Result struct {
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}
