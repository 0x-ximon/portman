package handlers

type Credentials struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
	OTP          string `json:"otp"`
}

type Payload struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}
