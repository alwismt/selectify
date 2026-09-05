package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type productFileRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

type ProductFileRepo interface {
	GetProductFileByProductID(ctx context.Context, proID uint) (productFile *model.ProductFile, err error)
	CreateProductFileWithTx(ctx context.Context, tx *sqlx.Tx, productFile *model.ProductFile) error
}

func NewProductFileRepo(db *db.DatabaseConnection) ProductFileRepo {
	return &productFileRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *productFileRepo) GetProductFileByProductID(ctx context.Context, proID uint) (*model.ProductFile, error) {
	var productFile model.ProductFile
	query := `SELECT * FROM product_file WHERE product_id = $1 AND is_primary = TRUE`

	err := r.roDb.GetContext(ctx, &productFile, query, proID)
	if err != nil {
		err = logger.Errorf(ctx, err, "error getting product file by proID %d", proID)
		return nil, err
	}

	return &productFile, nil
}

func (r *productFileRepo) CreateProductFileWithTx(ctx context.Context, tx *sqlx.Tx, productFile *model.ProductFile) error {
	query := `INSERT INTO product_file (product_file_id, product_id, variant_id, position, content_type, is_primary) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.ExecContext(ctx, query,
		productFile.FileID, productFile.ProductID, productFile.VariantID,
		productFile.Position, productFile.ContentType, productFile.IsPrimary)
	if err != nil {
		return logger.Errorf(ctx, err, "error creating product file")
	}
	return nil
}
