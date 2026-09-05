package repo

import (
	"context"
	"database/sql"
	"errors"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type productRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

type ProductRepo interface {
	GetProductByID(ctx context.Context, id uint) (*model.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*model.Product, error)
	GetProducts(ctx context.Context) (model.Products, error)
	GetProductsByIDs(ctx context.Context, ids []uint) (model.Products, error)

	GetProductsByMerchantID(ctx context.Context, id uint64) (model.Products, error)
	GetMerchantProductByID(ctx context.Context, merchantID uint64, slug uint) (*model.Product, error)
	CreateProductWithTx(ctx context.Context, tx *sqlx.Tx, product *model.Product) error

	GetProductsByCategoryID(ctx context.Context, categoryID uint) (model.Products, error)
	GetProductsByCategorySlug(ctx context.Context, categorySlug string) (model.Products, error)
}

func NewProductRepo(db *db.DatabaseConnection) ProductRepo {
	return &productRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (pr *productRepo) GetProductByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	query := `SELECT * FROM product WHERE product_id = $1 AND deleted_at IS NULL`

	err := pr.roDb.GetContext(ctx, &product, query, id)
	if err != nil {
		err = logger.Errorf(ctx, err, "error getting product by id %d", id)
		return nil, err
	}
	return &product, nil
}

func (pr *productRepo) GetMerchantProductByID(ctx context.Context, merchantID uint64, productID uint) (*model.Product, error) {
	var product model.Product
	query := `SELECT * FROM product WHERE product_id = $1 AND merchant_id = $2 AND deleted_at IS NULL`

	err := pr.roDb.GetContext(ctx, &product, query, productID, merchantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warnf(ctx, "error getting product by id %d, merchant id %d", productID, merchantID)
			return nil, err
		}
		return nil, logger.Errorf(ctx, err, "error getting product by id %d", productID)
	}
	return &product, nil
}

func (pr *productRepo) GetProductBySlug(ctx context.Context, slug string) (*model.Product, error) {
	var product model.Product
	query := `SELECT * FROM product WHERE slug = $1 AND deleted_at IS NULL`

	err := pr.roDb.GetContext(ctx, &product, query, slug)
	if err != nil {
		err = logger.Errorf(ctx, err, "error getting product by slug %s", slug)
		return nil, err
	}
	return &product, nil
}

func (pr *productRepo) GetProductsByIDs(ctx context.Context, ids []uint) (model.Products, error) {
	var products model.Products
	if len(ids) == 0 {
		return products, nil
	}
	query := `SELECT * FROM product WHERE product_id = ANY($1) AND deleted_at IS NULL`

	if err := pr.roDb.SelectContext(ctx, &products, query, pq.Array(ids)); err != nil {
		return products, logger.Errorf(ctx, err, "error getting products by id %v", ids)
	}
	return products, nil
}

type productRow struct {
	model.Product

	FileID          *uuid.UUID `db:"file_id"`
	FileContentType *string    `db:"content_type"`
	FilePosition    *uint      `db:"position"`
}

func (pr *productRepo) GetProducts(ctx context.Context) (model.Products, error) {
	const query = `SELECT
			p.*,
			pf.product_file_id AS file_id,
			pf.content_type AS content_type,
			pf.position AS position
		FROM product p
		LEFT JOIN product_file pf
			ON pf.product_id = p.product_id
			AND pf.variant_id IS NULL
			AND pf.is_primary = TRUE
		WHERE p.deleted_at IS NULL
		AND p.is_active = TRUE
		ORDER BY p.product_id DESC;`

	var rows []productRow
	if err := pr.roDb.SelectContext(ctx, &rows, query); err != nil {
		return nil, logger.Errorf(ctx, err, "error getting products")
	}

	products := make(model.Products, 0, len(rows))
	for _, row := range rows {
		product := row.Product

		if row.FileID != nil {
			product.ProductFile = &model.ProductFile{
				FileID:      *row.FileID,
				ProductID:   product.ID,
				VariantID:   nil,
				ContentType: *row.FileContentType,
				Position:    *row.FilePosition,
				IsPrimary:   true,
			}
		}
		products = append(products, product)
	}
	return products, nil
}

func (pr *productRepo) GetProductsByMerchantID(ctx context.Context, merchantID uint64) (model.Products, error) {
	const query = `SELECT
			p.*,
			pf.product_file_id AS file_id,
			pf.content_type AS content_type,
			pf.position AS position
		FROM product p
		LEFT JOIN product_file pf
			ON pf.product_id = p.product_id
			AND pf.variant_id IS NULL
			AND pf.is_primary = TRUE
		WHERE p.merchant_id = $1
		AND p.deleted_at IS NULL
		ORDER BY p.product_id DESC;`

	var rows []productRow
	if err := pr.roDb.SelectContext(ctx, &rows, query, merchantID); err != nil {
		return nil, logger.Errorf(ctx, err, "error getting products for merchant id %d", merchantID)
	}

	products := make(model.Products, 0, len(rows))
	for _, row := range rows {
		product := row.Product

		if row.FileID != nil {
			product.ProductFile = &model.ProductFile{
				FileID:      *row.FileID,
				ProductID:   product.ID,
				VariantID:   nil,
				ContentType: *row.FileContentType,
				Position:    *row.FilePosition,
				IsPrimary:   true,
			}
		}
		products = append(products, product)
	}
	return products, nil
}

func (pr *productRepo) CreateProductWithTx(ctx context.Context, tx *sqlx.Tx, product *model.Product) error {
	q := `INSERT INTO product (sku, slug, name, description, price_amount, merchant_id, is_active, in_stock)
		VALUES (:sku, :slug, :name, :description, :price_amount, :merchant_id, :is_active, :in_stock)
		RETURNING product_id, created_at, updated_at`

	q, args, err := sqlx.Named(q, product)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to bind product query")
	}

	q = tx.Rebind(q)

	if err = tx.QueryRowxContext(ctx, q, args...).StructScan(product); err != nil {
		return logger.Errorf(ctx, err, "failed to create product")
	}

	return nil
}

func (pr *productRepo) GetProductsByCategoryID(ctx context.Context, categoryID uint) (model.Products, error) {
	products := make(model.Products, 0)

	q := `SELECT
			p.product_id,
			p.sku,
			p.name,
			p.description,
			p.slug,
			p.price_amount,
			p.merchant_id,
			p.is_active,
			p.in_stock,
			p.created_at,
			p.updated_at
		FROM product p
		INNER JOIN product_category pc
			ON pc.product_id = p.product_id
		WHERE pc.category_id = $1
		  AND p.is_active = TRUE
		  AND p.deleted_at IS NULL
		ORDER BY p.created_at DESC`

	if err := pr.roDb.SelectContext(ctx, &products, q, categoryID); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get products by category")
	}

	return products, nil
}

func (pr *productRepo) GetProductsByCategorySlug(ctx context.Context, categorySlug string) (model.Products, error) {
	products := make(model.Products, 0)

	q := `SELECT
			p.product_id,
			p.sku,
			p.name,
			p.description,
			p.slug,
			p.price_amount,
			p.merchant_id,
			p.is_active,
			p.in_stock,
			p.created_at,
			p.updated_at
		FROM product p
		INNER JOIN product_category pc
			ON pc.product_id = p.product_id
		INNER JOIN category c
			ON c.category_id = pc.category_id
		WHERE c.slug = $1
		  AND c.is_active = TRUE
		  AND p.is_active = TRUE
		  AND p.deleted_at IS NULL
		ORDER BY p.created_at DESC`

	if err := pr.roDb.SelectContext(ctx, &products, q, categorySlug); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get products by category slug")
	}

	return products, nil
}
