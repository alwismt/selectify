package httpx

import (
	"errors"
	"net/http"
)

func SendInvalidData(w http.ResponseWriter) {
	SendError(w, errors.New("invalid error"))
}

//func SendJson(w http.ResponseWriter, status int, d interface{}) error {
//	return sendable.SendJson(w, status, d)
//}
