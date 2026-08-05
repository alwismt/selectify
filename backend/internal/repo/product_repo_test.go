package repo_test

import (
	"fmt"
	"math/rand"
	"testing"

	"alwis.dev/selectify/internal/model"

	"github.com/stretchr/testify/require"
)

func testCreateProduct(t *testing.T) *model.Product {
	product := new(model.Product)
	description := "Test product description"
	slug := fmt.Sprintf("test-product-slug-%d", rand.Int())
	sku := fmt.Sprintf("TSET-SKU-%d", rand.Int())

	insertQuery := `INSERT INTO product (sku, name, description, slug, price_amount, currency, is_active, in_stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	err := dbConn.QueryRowx(insertQuery, sku, "Test Product", description, slug, 99.99, "EUR", true, true).StructScan(product)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = dbConn.Exec("DELETE FROM product WHERE product_id = $1", product.ID)
		require.NoError(t, err)
	})
	return product
}

func TestProductRepo_GetProductByID(t *testing.T) {
	expectedProduct := testCreateProduct(t)
	product, err := productRepo.GetProductByID(ctx, expectedProduct.ID)

	require.NoError(t, err)
	require.NotNil(t, product)
	require.Equal(t, expectedProduct.ID, product.ID)
	require.Equal(t, expectedProduct.SKU, product.SKU)
	require.Equal(t, expectedProduct.Name, product.Name)
	require.NotNil(t, product.Description)
	require.Equal(t, *expectedProduct.Description, *product.Description)
	require.NotNil(t, product.Slug)
	require.Equal(t, *expectedProduct.Slug, *product.Slug)
	require.Equal(t, expectedProduct.PriceAmount, product.PriceAmount)
	require.Equal(t, expectedProduct.Currency, product.Currency)
	require.Equal(t, expectedProduct.IsActive, product.IsActive)
	require.Equal(t, expectedProduct.InStock, product.InStock)

	product, err = productRepo.GetProductByID(ctx, 999999)
	require.Error(t, err, "GetProductByID should return an error for non-existent product")
	require.Nil(t, product, "Product should be nil when not found")
}

func TestProductRepo_GetProducts(t *testing.T) {
	products, err := productRepo.GetProducts(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, products)
	require.Len(t, products, 2)

	require.Equal(t, testProduct.ID, products[0].ID)
	require.Equal(t, testProduct.Name, products[0].Name)
	require.Equal(t, testProduct.SKU, products[0].SKU)
	require.Equal(t, testProduct2.ID, products[1].ID)
	require.Equal(t, testProduct2.Name, products[1].Name)
	require.Equal(t, testProduct2.SKU, products[1].SKU)
}
