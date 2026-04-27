package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	codeBadRequest      = "BAD_REQUEST"
	codeUnauthorized    = "UNAUTHORIZED"
	codeForbidden       = "FORBIDDEN"
	codeNotFound        = "NOT_FOUND"
	codeInternal        = "INTERNAL_SERVER_ERROR"
	codeNotImplemented  = "NOT_IMPLEMENTED"
	codeInvalidProvider = "INVALID_PROVIDER"
)

func SendSuccess(w http.ResponseWriter, data Data) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Payload{
		Status: "success",
		Data:   &data,
	})
}

func SendFailure(w http.ResponseWriter, status int, code string, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Payload{
		Status: "failure",
		Error: &Error{
			Code:   code,
			Detail: detail,
		},
	})
}
