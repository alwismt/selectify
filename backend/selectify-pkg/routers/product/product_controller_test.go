package product_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/testkit/test"
)

func testProductEqual(t *testing.T, expected, actual model.Product) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.SKU, actual.SKU)
	require.Equal(t, expected.Name, actual.Name)
	if expected.Description != nil {
		require.NotNil(t, actual.Description)
		require.Equal(t, *expected.Description, *actual.Description)
	}
	if expected.Slug != nil {
		require.NotNil(t, actual.Slug)
		require.Equal(t, *expected.Slug, *actual.Slug)
	}
	require.Equal(t, expected.PriceAmount, actual.PriceAmount)
	require.Equal(t, expected.Currency, actual.Currency)
	require.Equal(t, expected.IsActive, actual.IsActive)
	require.Equal(t, expected.InStock, actual.InStock)
}

func TestController_GetProducts(t *testing.T) {
	getProductsPath := "/api/v1/products"
	expectedProducts := new(model.Products)
	err := ts.DB.RwDb.Select(expectedProducts, `SELECT * FROM product`)
	require.NoError(t, err)

	products := new(model.Products)

	resp := test.DoGet(t, ts.S.URL, getProductsPath, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	err = test.BodyAsJSON(resp, products)
	require.NoError(t, err)
	require.NotNil(t, products)
	require.NotEmpty(t, products)
	require.Len(t, *products, 2)
	require.Len(t, *expectedProducts, 2)

	testProductEqual(t, (*expectedProducts)[0], (*products)[0])
	testProductEqual(t, (*expectedProducts)[1], (*products)[1])
}

func TestController_GetProductById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedProducts := new(model.Products)
		err := ts.DB.RwDb.Select(expectedProducts, `SELECT * FROM product LIMIT 1`)
		require.NoError(t, err)
		require.NotEmpty(t, *expectedProducts)

		expectedProduct := (*expectedProducts)[0]
		getProductByIdPath := fmt.Sprintf("/api/v1/products/%d", expectedProduct.ID)

		product := new(model.Product)

		resp := test.DoGet(t, ts.S.URL, getProductByIdPath, nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		err = test.BodyAsJSON(resp, product)
		require.NoError(t, err)
		require.NotNil(t, product)

		testProductEqual(t, expectedProduct, *product)
	})

	t.Run("NotFound", func(t *testing.T) {
		getProductByIdPath := "/api/v1/products/999999"

		resp := test.DoGet(t, ts.S.URL, getProductByIdPath, nil)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		getProductByIdPath := "/api/v1/products/invalid"

		resp := test.DoGet(t, ts.S.URL, getProductByIdPath, nil)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
