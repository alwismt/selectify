package user

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	controller_utils "alwis.dev/selectify/internal/httpx/controller"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/selectify-pkg/app"
)

const maxProfileImageSize = 5 << 20

func (c *controller) UserInfo(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	s.User.UserRole = s.UserRole
	if err := httpx.SendJson(w, http.StatusOK, s.User); err != nil {
		_ = logger.Error(r.Context(), err, "failed to send user info response")
	}
}

func (c *controller) GetUserImage(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
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

func (c *controller) UpdateUserImage(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	file, err := controller_utils.GetMultiPartFile(r, "image", maxProfileImageSize, "image")
	if err != nil {
		_ = logger.Error(r.Context(), err, "image is required")
		httpx.SendError(w, err)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			_ = logger.Error(r.Context(), err, "Failed to close file")
		}
	}()

	userFile, err := c.userService.UpsertUserImage(r.Context(), s.User, file)
	if err != nil {
		_ = logger.Error(r.Context(), err, "failed to update user image")
		httpx.SendError(w, err)
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, userFile); err != nil {
		_ = logger.Error(r.Context(), err, "failed to send user file response")
	}
}

func (c *controller) DeleteUserImage(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	if err := c.userService.DeleteUserImage(r.Context(), s.User); err != nil {
		_ = logger.Error(r.Context(), err, "failed to delete user image")
		httpx.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetDefaultAddress(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
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
