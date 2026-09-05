package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/lib/pq"

	"alwis.dev/selectify/internal/helper/validation"
	"alwis.dev/selectify/internal/httpx"
	controller_utils "alwis.dev/selectify/internal/httpx/controller"
	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

func (c *controller) GetProducts(w http.ResponseWriter, r *http.Request, _ *model.LoggedInSession, m *model.Merchant) {
	ctx := r.Context()
	products, err := c.productService.GetProductsByMerchant(ctx, m)
	if err != nil {
		_ = logger.Error(ctx, err, "failed to get products")
		httpx.SendError(w, httpx.ErrBadRequest)
		return
	}

	if err = httpx.NewJsonSender(products, http.StatusOK).Send(w); err != nil {
		_ = logger.Error(ctx, err, "failed to send response")
		httpx.SendError(w, httpx.ErrBadRequest)
	}

	return
}

func (c *controller) GetProductByID(w http.ResponseWriter, r *http.Request, _ *model.LoggedInSession, _ *model.Merchant, p *model.Product) {
	productFile, err := c.productFileService.GetProductFileByProductID(r.Context(), p.ID)
	if err == nil && productFile != nil {
		p.ProductFile = productFile
	}

	if err = httpx.NewJsonSender(p, http.StatusOK).Send(w); err != nil {
		_ = logger.Error(r.Context(), err, "failed to send response")
		httpx.SendError(w, httpx.ErrInternalServer)
	}

	return
}

func (c *controller) CreateProduct(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession, m *model.Merchant) {
	req := new(request.ProductRequest)
	ctx := r.Context()
	fieldErrors := make(map[string]string)

	req.Name = r.FormValue("name")
	req.Description = r.FormValue("description")
	req.SKU = r.FormValue("sku")
	categoryValues := r.Form["categoryIds"]
	priceValue := r.FormValue("price")

	req.Sanitize()

	for _, value := range categoryValues {
		categoryID, err := strconv.ParseUint(value, 10, strconv.IntSize)
		if err != nil {
			fieldErrors["categoryIds"] = "Invalid value"
			break
		}
		req.CategoryIDs = append(req.CategoryIDs, uint(categoryID))
	}

	if priceValue != "" {
		price, err := strconv.ParseUint(priceValue, 10, strconv.IntSize)
		if err != nil {
			_ = logger.Error(ctx, err, "Invalid price")
			fieldErrors["price"] = "Invalid value"
		} else {
			req.Price = uint(price)
		}
	}

	var err error
	if err = validation.Validate.Struct(req); err != nil {
		_ = logger.Errorf(ctx, err, "Validation error")
		for field, message := range httpx.ValidationErrorsToMap(err) {
			if _, exists := fieldErrors[field]; exists {
				continue
			}
			fieldErrors[field] = message
		}

	}

	if len(fieldErrors) > 0 {
		_ = httpx.SendJson(w, http.StatusBadRequest, fieldErrors)
		return
	}

	req.Image, err = controller_utils.GetMultiPartFile(r, "image", controller_utils.MaxProfileImageSize, "image")
	if err != nil {
		_ = logger.Error(r.Context(), err, "image is required")
		fieldErrors["image"] = "Image is required"
		_ = httpx.SendJson(w, http.StatusBadRequest, fieldErrors)
		return
	}

	defer func() {
		err = req.Image.Close()
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

	pro, err := c.productService.CreateProduct(ctx, req, m)
	if err != nil {
		_ = logger.Error(ctx, err, "Failed to create product")
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			_ = httpx.SendJson(w, http.StatusBadRequest, map[string]string{
				"sku": "SKU already exists",
			})
			return
		}
		httpx.SendError(w, httpx.ErrBadRequest)
		return
	}

	if err = httpx.SendJson(w, http.StatusCreated, pro); err != nil {
		_ = logger.Error(r.Context(), err, "Failed to send response")
	}

	return
}
