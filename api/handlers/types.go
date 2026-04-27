package handlers

type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type Data = any

type Payload struct {
	Status string `json:"status"`
	Error  *Error `json:"error,omitempty"`
	Data   *Data  `json:"data,omitempty"`
}
