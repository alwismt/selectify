package httpx

import (
	"alwis.dev/selectify/internal/logger"
	"encoding/json"
	"io"
	"net/http"
)

func MustDecodeJson(w http.ResponseWriter, r *http.Request, i interface{}) error {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = logger.Error(ctx, err, "failed to read body")
		SendInvalidData(w)
		return err
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			_ = logger.Error(ctx, err, "failed to close body")
		}
	}()

	err = json.Unmarshal(body, i)
	if err != nil {
		_ = logger.Error(ctx, err, "failed to unmarshal body")
		SendInvalidData(w)
		return err
	}

	return nil
}
