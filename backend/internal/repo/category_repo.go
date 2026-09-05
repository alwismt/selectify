package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type categoryRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

type CategoryRepo interface {
	AddProductCategoryWithTx(ctx context.Context, tx *sqlx.Tx, productID, categoryID uint) error
	AddProductCategoriesWithTx(ctx context.Context, tx *sqlx.Tx, productID uint, categoryIDs []uint) error

	GetCategoryByID(ctx context.Context, categoryID uint) (*model.Category, error)
	GetCategories(ctx context.Context) (model.Categories, error)
}

func NewCategoryRepo(db *db.DatabaseConnection) CategoryRepo {
	return &categoryRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *categoryRepo) AddProductCategoryWithTx(ctx context.Context, tx *sqlx.Tx, productID, categoryID uint) error {
	q := `INSERT INTO product_category (product_id, category_id)
		VALUES ($1, $2)`

	if _, err := tx.ExecContext(ctx, q, productID, categoryID); err != nil {
		return logger.Errorf(ctx, err, "failed to add product category")
	}

	return nil
}

func (r *categoryRepo) AddProductCategoriesWithTx(ctx context.Context, tx *sqlx.Tx, productID uint, categoryIDs []uint) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	args := make([]any, 0, len(categoryIDs)+1)
	args = append(args, productID)

	values := make([]string, 0, len(categoryIDs))
	for i, categoryID := range categoryIDs {
		values = append(values, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, categoryID)
	}

	q := `INSERT INTO product_category (product_id, category_id)
		VALUES` + strings.Join(values, ", ")

	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return logger.Errorf(ctx, err, "failed to add product categories")
	}

	return nil
}

func (r *categoryRepo) GetCategoryByID(ctx context.Context, categoryID uint) (*model.Category, error) {
	category := new(model.Category)
	q := `SELECT category_id, name, slug, parent_id, is_active, created_at, updated_at
		FROM category
		WHERE category_id = $1`

	if err := r.roDb.GetContext(ctx, category, q, categoryID); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get category")
	}

	return category, nil
}

func (r *categoryRepo) GetCategories(ctx context.Context) (model.Categories, error) {
	categories := make(model.Categories, 0)

	q := `SELECT category_id, name, slug, parent_id, is_active, created_at, updated_at
		FROM category
		WHERE is_active = TRUE
		ORDER BY name`

	if err := r.roDb.SelectContext(ctx, &categories, q); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get categories")
	}

	return categories, nil
}
