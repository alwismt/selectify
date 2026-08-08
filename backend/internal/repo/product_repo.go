package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"

	"github.com/jmoiron/sqlx"
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
}

func NewProductRepo(db *db.DatabaseConnection) ProductRepo {
	return &productRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (pr *productRepo) GetProductByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	query := `SELECT * FROM product WHERE product_id = $1`

	err := pr.roDb.GetContext(ctx, &product, query, id)
	if err != nil {
		err = logger.Errorf(ctx, err, "error getting product by id %d", id)
		return nil, err
	}
	return &product, nil
}

func (pr *productRepo) GetProductBySlug(ctx context.Context, slug string) (*model.Product, error) {
	var product model.Product
	query := `SELECT * FROM product WHERE slug = $1`

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
	query := `SELECT * FROM product WHERE product_id = ANY($1)`

	if err := pr.roDb.SelectContext(ctx, &products, query, pq.Array(ids)); err != nil {
		return products, logger.Errorf(ctx, err, "error getting products by id %v", ids)
	}
	return products, nil
}

func (pr *productRepo) GetProducts3(ctx context.Context) (model.Products, error) {
	var products model.Products

	query := `SELECT * FROM product ORDER BY product_id DESC`
	err := pr.roDb.SelectContext(ctx, &products, query)
	if err != nil {
		err = logger.Errorf(ctx, err, "error getting products")
		return nil, err
	}
	return products, nil
}

type productRow struct {
	model.Product

	FileID          *uuid.UUID `db:"file_id"`
	FileContentType *string    `db:"content_type"`
	FilePosition    *uint      `db:"position"`
}

func (pr *productRepo) GetProducts(
	ctx context.Context,
) (model.Products, error) {
	const query = `
		SELECT
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
		ORDER BY p.product_id DESC;
	`

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
