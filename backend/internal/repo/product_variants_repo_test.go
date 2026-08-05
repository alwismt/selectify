package repo_test

import (
	"fmt"
	"math/rand"
	"testing"

	"alwis.dev/selectify/internal/model"

	"github.com/stretchr/testify/require"
)

func testCreateProductVariant(t *testing.T, productID uint) *model.ProductVariant {
	variant := new(model.ProductVariant)
	sku := fmt.Sprintf("TEST-VARIANT-SKU-%d", rand.Int())
	priceAmount := 29.99

	insertQuery := `INSERT INTO product_variants (product_id, sku, price_amount, currency, is_active)
		VALUES ($1, $2, $3, $4, $5) RETURNING *`

	err := dbConn.QueryRowx(insertQuery, productID, sku, priceAmount, "EUR", true).StructScan(variant)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = dbConn.Exec("DELETE FROM product_variants WHERE id = $1", variant.ID)
		require.NoError(t, err)
	})
	return variant
}

func testCreateProductVariantAttribute(t *testing.T, variantID uint, name, value string) *model.ProductVariantAttribute {
	attr := new(model.ProductVariantAttribute)

	insertQuery := `INSERT INTO product_variant_attributes (variant_id, name, value)
		VALUES ($1, $2, $3) RETURNING *`

	err := dbConn.QueryRowx(insertQuery, variantID, name, value).StructScan(attr)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = dbConn.Exec("DELETE FROM product_variant_attributes WHERE id = $1", attr.ID)
		require.NoError(t, err)
	})
	return attr
}

func testCreateInventory(t *testing.T, variantID uint, stockQty, reservedQty uint) {
	insertQuery := `INSERT INTO inventory (variant_id, stock_qty, reserved_qty)
		VALUES ($1, $2, $3)
		ON CONFLICT (variant_id) DO UPDATE SET stock_qty = $2, reserved_qty = $3`

	_, err := dbConn.Exec(insertQuery, variantID, stockQty, reservedQty)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = dbConn.Exec("DELETE FROM inventory WHERE variant_id = $1", variantID)
		require.NoError(t, err)
	})
}

func TestProductVariantsRepo_GetVariantsForProduct(t *testing.T) {
	// Create a product first
	product := testCreateProduct(t)

	// Create first variant with attributes and inventory
	variant1 := testCreateProductVariant(t, product.ID)
	testCreateProductVariantAttribute(t, variant1.ID, "Color", "Red")
	testCreateProductVariantAttribute(t, variant1.ID, "Size", "Large")
	testCreateInventory(t, variant1.ID, 100, 10)

	// Create second variant with attributes but no inventory
	variant2 := testCreateProductVariant(t, product.ID)
	testCreateProductVariantAttribute(t, variant2.ID, "Color", "Blue")
	testCreateProductVariantAttribute(t, variant2.ID, "Size", "Medium")

	// Create third variant without attributes or inventory
	variant3 := testCreateProductVariant(t, product.ID)

	// Test: Get variants for the product
	variants, err := productVariantsRepo.GetVariantsForProduct(ctx, product.ID)

	require.NoError(t, err)
	require.NotNil(t, variants)
	require.Len(t, variants, 3)

	// Verify first variant with attributes and inventory
	foundVariant1 := findVariantByID(variants, variant1.ID)
	require.NotNil(t, foundVariant1)
	require.Equal(t, variant1.ID, foundVariant1.ID)
	require.Equal(t, variant1.ProductID, foundVariant1.ProductID)
	require.Equal(t, variant1.SKU, foundVariant1.SKU)
	require.NotNil(t, foundVariant1.PriceAmount)
	require.Equal(t, *variant1.PriceAmount, *foundVariant1.PriceAmount)
	require.Equal(t, variant1.Currency, foundVariant1.Currency)
	require.True(t, foundVariant1.IsActive)
	require.Len(t, foundVariant1.ProductVariantAttributes, 2)
	require.NotNil(t, foundVariant1.Inventory)
	require.Equal(t, uint(100), foundVariant1.Inventory.StockQty)
	require.Equal(t, uint(10), foundVariant1.Inventory.ReservedQty)

	// Verify attributes are correctly parsed
	attrMap := make(map[string]string)
	for _, attr := range foundVariant1.ProductVariantAttributes {
		attrMap[attr.Name] = attr.Value
	}
	require.Equal(t, "Red", attrMap["Color"])
	require.Equal(t, "Large", attrMap["Size"])

	// Verify second variant with attributes but no inventory
	foundVariant2 := findVariantByID(variants, variant2.ID)
	require.NotNil(t, foundVariant2)
	require.Equal(t, variant2.ID, foundVariant2.ID)
	require.Len(t, foundVariant2.ProductVariantAttributes, 2)
	require.Empty(t, foundVariant2.Inventory)

	// Verify third variant without attributes or inventory
	foundVariant3 := findVariantByID(variants, variant3.ID)
	require.NotNil(t, foundVariant3)
	require.Equal(t, variant3.ID, foundVariant3.ID)
	require.Len(t, foundVariant3.ProductVariantAttributes, 0)
	require.Empty(t, foundVariant3.Inventory)
}

func TestProductVariantsRepo_GetVariantsForProduct_NoVariants(t *testing.T) {
	// Create a product without variants
	product := testCreateProduct(t)

	variants, err := productVariantsRepo.GetVariantsForProduct(ctx, product.ID)

	require.NoError(t, err)
	require.NotNil(t, variants)
	require.Len(t, variants, 0)
}

func TestProductVariantsRepo_GetVariantsForProduct_ExcludesDeleted(t *testing.T) {
	product := testCreateProduct(t)

	// Create an active variant
	activeVariant := testCreateProductVariant(t, product.ID)

	// Create a deleted variant
	deletedVariant := testCreateProductVariant(t, product.ID)
	_, err := dbConn.Exec("UPDATE product_variants SET deleted_at = NOW() WHERE id = $1", deletedVariant.ID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = dbConn.Exec("DELETE FROM product_variants WHERE id = $1", deletedVariant.ID)
	})

	variants, err := productVariantsRepo.GetVariantsForProduct(ctx, product.ID)

	require.NoError(t, err)
	require.NotNil(t, variants)
	require.Len(t, variants, 1)
	require.Equal(t, activeVariant.ID, variants[0].ID)
}

func TestProductVariantsRepo_GetVariantsForProduct_ExcludesInactive(t *testing.T) {
	product := testCreateProduct(t)

	// Create an active variant
	activeVariant := testCreateProductVariant(t, product.ID)

	// Create an inactive variant
	inactiveVariant := testCreateProductVariant(t, product.ID)
	_, err := dbConn.Exec("UPDATE product_variants SET is_active = FALSE WHERE id = $1", inactiveVariant.ID)
	require.NoError(t, err)

	variants, err := productVariantsRepo.GetVariantsForProduct(ctx, product.ID)

	require.NoError(t, err)
	require.NotNil(t, variants)
	require.Len(t, variants, 1)
	require.Equal(t, activeVariant.ID, variants[0].ID)
}

func TestProductVariantsRepo_GetVariantsForProduct_NonExistentProduct(t *testing.T) {
	variants, err := productVariantsRepo.GetVariantsForProduct(ctx, 999999)

	require.NoError(t, err)
	require.NotNil(t, variants)
	require.Len(t, variants, 0)
}

func findVariantByID(variants model.ProductVariants, id uint) *model.ProductVariant {
	for i := range variants {
		if variants[i].ID == id {
			return &variants[i]
		}
	}
	return nil
}
