package product_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/testkit/test"
)

func TestController_GetVariantsForProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expectedProducts := new(model.Products)
		err := ts.DB.RwDb.Select(expectedProducts, `SELECT * FROM product LIMIT 1`)
		require.NoError(t, err)
		require.NotEmpty(t, *expectedProducts)

		expectedProduct := (*expectedProducts)[0]
		getVariantsPath := fmt.Sprintf("/api/v1/products/%d/variants", expectedProduct.ID)

		variants := new(model.ProductVariants)

		resp := test.DoGet(t, ts.S.URL, getVariantsPath, nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		err = test.BodyAsJSON(resp, variants)
		require.NoError(t, err)
		require.NotNil(t, variants)

		var expectedCount int
		err = ts.DB.RwDb.Get(&expectedCount, `
			SELECT COUNT(*) FROM product_variants 
			WHERE product_id = $1 AND deleted_at IS NULL AND is_active = TRUE`, expectedProduct.ID)
		require.NoError(t, err)

		require.Len(t, *variants, expectedCount)

		if expectedCount > 0 {
			for _, variant := range *variants {
				require.Equal(t, expectedProduct.ID, variant.ProductID)
				require.NotEmpty(t, variant.SKU)
				require.NotZero(t, variant.ID)
				require.True(t, variant.IsActive)
			}
		}
	})

	t.Run("NoVariants", func(t *testing.T) {
		expectedProducts := new(model.Products)
		err := ts.DB.RwDb.Select(expectedProducts, `SELECT * FROM product WHERE product_id NOT IN (
			SELECT DISTINCT product_id FROM product_variants WHERE deleted_at IS NULL AND is_active = TRUE
		) LIMIT 1`)
		if err != nil || len(*expectedProducts) == 0 {
			t.Skip("No product without variants found")
			return
		}

		expectedProduct := (*expectedProducts)[0]
		getVariantsPath := fmt.Sprintf("/api/v1/products/%d/variants", expectedProduct.ID)

		variants := new(model.ProductVariants)

		resp := test.DoGet(t, ts.S.URL, getVariantsPath, nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		err = test.BodyAsJSON(resp, variants)
		require.NoError(t, err)
		require.NotNil(t, variants)
		require.Len(t, *variants, 0)
	})

	t.Run("NotFound", func(t *testing.T) {
		getVariantsPath := "/api/v1/products/999999/variants"

		resp := test.DoGet(t, ts.S.URL, getVariantsPath, nil)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		getVariantsPath := "/api/v1/products/invalid/variants"

		resp := test.DoGet(t, ts.S.URL, getVariantsPath, nil)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
