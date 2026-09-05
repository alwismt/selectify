package controller_utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

var (
	ErrMissingFile        = errors.New("FILE_MISSING")
	ErrInvalidContentType = errors.New("FILE_CONTENT_TYPE_NOT_ALLOWED")
	ErrFileSizeLimit      = errors.New("FILE_EXCEEDS_SIZE_LIMIT")
)

const MaxProfileImageSize = 5 << 20

func GetMultiPartFile(r *http.Request, key string, maxSize int64, requiredContentType string) (model.FileStream, error) {
	if err := r.ParseMultipartForm(256 << 10); err != nil {
		return nil, logger.Errorf(r.Context(), err, "failed to parse multipart form")
	}

	file, header, err := r.FormFile(key)
	if err != nil {
		return nil, fmt.Errorf("key %q missing: %w", key, ErrMissingFile)
	}

	if header.Size == 0 || (maxSize != -1 && header.Size > maxSize) {
		_ = file.Close()
		return nil, logger.Errorf(r.Context(), nil, "file size (%d bytes) exceeds allowed size: %v", header.Size, ErrFileSizeLimit)
	}

	buf := make([]byte, 512)
	n, readErr := file.Read(buf)
	if readErr != nil && readErr != io.EOF {
		_ = file.Close()
		return nil, logger.Errorf(r.Context(), nil, "failed to read file: %v", readErr)
	}

	contentType := http.DetectContentType(buf[:n])

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		_ = file.Close()
		return nil, logger.Errorf(r.Context(), nil, "failed to seek file: %v", err)
	}

	switch requiredContentType {
	case "":
	case "image":
		switch contentType {
		case "image/jpeg", "image/png", "image/webp":
		default:
			_ = file.Close()
			return nil, logger.Error(r.Context(), nil, ErrInvalidContentType.Error())
		}
	case "pdf":
		if contentType != "application/pdf" {
			_ = file.Close()
			return nil, logger.Error(r.Context(), nil, ErrInvalidContentType.Error())
		}
	default:
		_ = file.Close()
		return nil, logger.Error(r.Context(), nil, ErrInvalidContentType.Error())
	}

	return model.NewFileStream(file, contentType, header.Size), nil
}
