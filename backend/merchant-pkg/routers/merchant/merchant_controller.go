package merchant

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"alwis.dev/selectify/internal/httpx"
	controller_utils "alwis.dev/selectify/internal/httpx/controller"
	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

const maxProfileImageSize = 5 << 20

func (c *controller) GetMerchant(w http.ResponseWriter, r *http.Request, _ *model.LoggedInSession, merchant *model.Merchant) {
	if err := httpx.SendJson(w, http.StatusOK, merchant); err != nil {
		_ = logger.Error(r.Context(), err, "Failed to send response")
	}
	return
}

func (c *controller) UpdateMerchant(w http.ResponseWriter, r *http.Request, _ *model.LoggedInSession, merchant *model.Merchant) {
	req := new(request.UpdateMerchantRequest)
	ctx := r.Context()

	if err := httpx.MustDecodeJson(w, r, req); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to decode body")
		return
	}

	if req.Empty() {
		_ = logger.Error(ctx, fmt.Errorf("request is nil"), "Failed to decode body")
		httpx.SendError(w, httpx.ErrInvalidRequest)
		return
	}
	req.Sanitize()

	if err := validator.New().Struct(req); err != nil {
		_ = logger.Errorf(ctx, err, "Validation error")
		fieldErrors := httpx.ValidationErrorsToMap(err)
		_ = httpx.SendJson(w, http.StatusBadRequest, fieldErrors)
		return
	}

	if err := c.merchantService.UpdateMerchant(ctx, merchant, req); err != nil {
		if errors.Is(err, httpx.ErrInvalidCountry) {
			httpx.SendError(w, httpx.ErrInvalidCountry)
			return
		}
		_ = logger.Errorf(ctx, err, "Failed to update merchant")
		httpx.SendError(w, httpx.ErrBadRequest)
		return
	}

	if err := httpx.SendJson(w, http.StatusOK, merchant); err != nil {
		_ = logger.Error(r.Context(), err, "Failed to send response")
	}
	return
}

func (c *controller) UpdateMerchantLogo(w http.ResponseWriter, r *http.Request, _ *model.LoggedInSession, merchant *model.Merchant) {
	file, err := controller_utils.GetMultiPartFile(r, "image", maxProfileImageSize, "image")
	if err != nil {
		_ = logger.Error(r.Context(), err, "image is required")
		httpx.SendError(w, err)
		return
	}
	ctx := r.Context()
	defer func() {
		err = file.Close()
		if err != nil {
			_ = logger.Error(ctx, err, "Failed to close file")
		}
	}()

	if r.MultipartForm != nil {
		defer func() {
			if err = r.MultipartForm.RemoveAll(); err != nil {
				_ = logger.Error(ctx, err, "Failed to remove multipart form")
			}
		}()
	}
	if err = c.merchantService.UpdateMerchantImage(ctx, merchant, file); err != nil {
		_ = logger.Error(ctx, err, "Failed to update merchant image")
		httpx.SendError(w, httpx.ErrBadRequest)
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, merchant.Logo); err != nil {
		_ = logger.Error(ctx, err, "Failed to send response")
	}
	return
}

func (c *controller) GetMerchantCountries(w http.ResponseWriter, r *http.Request, _ *model.LoggedInSession) {
	countries, err := c.countryService.GetCountries(r.Context())
	if err != nil {
		_ = logger.Error(r.Context(), err, "Failed to get countries")
		httpx.SendError(w, err)
		return
	}
	if err = httpx.SendJson(w, http.StatusOK, countries); err != nil {
		_ = logger.Error(r.Context(), err, "Failed to send response")
	}
	return
}
