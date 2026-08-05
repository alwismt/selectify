package httpx

import (
	"context"
	"encoding/json"
	"net/http"

	"alwis.dev/selectify/internal/logger"
)

const ContentTypeApplicationJson = "application/json"

type jsonSender struct {
	data   interface{}
	status int
}

type Sender interface {
	Send(w http.ResponseWriter) error
}

func NewJsonSender(data interface{}, status int) Sender {
	return &jsonSender{
		data:   data,
		status: status,
	}
}

func SendJson(w http.ResponseWriter, status int, data interface{}) error {
	return NewJsonSender(data, status).Send(w)
}

func (s *jsonSender) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", ContentTypeApplicationJson)
	w.WriteHeader(s.status)

	e := json.NewEncoder(w)
	err := e.Encode(s.data)
	if err != nil {
		return err
	}

	return nil
}

func SendStatus(w http.ResponseWriter, statusCode int, status string) {
	if err := SendJson(w, statusCode, StatusResponse{Status: "ok", Message: status}); err != nil {
		_ = logger.Error(context.Background(), err, "status send failed")
	}
}

func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
