package user

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/selectify-pkg/app"
)

const maxProfileImageSize = 5 << 20

func (c *controller) UserInfo(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	s.User.UserRole = s.UserRole
	if err := httpx.SendJson(w, http.StatusOK, s.User); err != nil {
		_ = logger.Error(r.Context(), err, "failed to send user info response")
	}
}

func (c *controller) GetUserImage(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	userFile, err := c.userService.GetUserImage(r.Context(), s.User)
	if err != nil {
		_ = logger.Error(r.Context(), err, "failed to get user image")
		httpx.SendError(w, err)
		return
	}

	if userFile == nil {
		if err = httpx.SendJson(w, http.StatusOK, map[string]any{}); err != nil {
			_ = logger.Error(r.Context(), err, "failed to send empty user file response")
		}
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, userFile); err != nil {
		_ = logger.Error(r.Context(), err, "failed to send user file response")
	}
}

func (c *controller) UpdateUserImage(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	r.Body = http.MaxBytesReader(w, r.Body, maxProfileImageSize)

	file, _, err := r.FormFile("image")

	if err != nil {
		_ = logger.Error(r.Context(), err, "image is required")
		http.Error(w, "image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()
	contentType, err := detectImageContentType(file)
	if err != nil {
		_ = logger.Error(r.Context(), err, "invalid image")
		http.Error(w, "invalid image", http.StatusBadRequest)
		return
	}

	userFile, err := c.userService.UpsertUserImage(r.Context(), s.User, file, contentType)
	if err != nil {
		_ = logger.Error(r.Context(), err, "failed to update user image")
		httpx.SendError(w, err)
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, userFile); err != nil {
		_ = logger.Error(r.Context(), err, "failed to send user file response")
	}
}

func (c *controller) DeleteUserImage(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	if err := c.userService.DeleteUserImage(r.Context(), s.User); err != nil {
		_ = logger.Error(r.Context(), err, "failed to delete user image")
		httpx.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func detectImageContentType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read image header: %w", err)
	}

	contentType := http.DetectContentType(buffer[:n])

	switch contentType {
	case "image/jpeg", "image/jpg", "image/png", "image/webp":
	default:
		return "", fmt.Errorf("unsupported image type: %s", contentType)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("reset image reader: %w", err)
	}

	return contentType, nil
}

func GetDefaultAddress(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	ctx := r.Context()

	addr, err := app.Repository().UserAddressRepo.GetDefaultByUserID(ctx, s.User.ID)
	if err != nil {
		_ = logger.Error(ctx, err, "failed to get default user address")
		httpx.SendError(w, err)
		return
	}

	if addr == nil {
		if err = httpx.SendJson(w, http.StatusOK, map[string]any{}); err != nil {
			_ = logger.Error(ctx, err, "failed to send empty address response")
		}
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, addr); err != nil {
		_ = logger.Error(ctx, err, "failed to send user address response")
	}
}
