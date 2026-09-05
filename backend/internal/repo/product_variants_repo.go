package repo

import (
	"context"
	"encoding/json"

	"github.com/lib/pq"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"

	"github.com/jmoiron/sqlx"
)

type productVariantsRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

type ProductVariantsRepo interface {
	GetVariantsForProduct(ctx context.Context, productID uint) (model.ProductVariants, error)
	GetVariantsByID(ctx context.Context, variantId uint) (*model.ProductVariant, error)
	GetVariantsByIDs(ctx context.Context, variantIDs []uint) (model.ProductVariants, error)
	ReserveStockForCheckout(ctx context.Context, tx *sqlx.Tx, variantID uint, qty uint) error
	GetAttrByVariantsIDs(ctx context.Context, variantIDs []uint) (model.ProductVariantAttributes, error)
}

func NewProductVariantRepo(db *db.DatabaseConnection) ProductVariantsRepo {
	return &productVariantsRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

type variantRow struct {
	model.ProductVariant `db:",inline"`

	AttrsJSON []byte          `db:"product_variant_attributes"`
	FilesJSON json.RawMessage `db:"product_files"`
}

func (r *productVariantsRepo) GetVariantsForProduct(ctx context.Context, productID uint) (model.ProductVariants, error) {
	const q = `SELECT
			v.id,
			v.product_id,
			v.sku,
			v.price_amount,
			v.created_at,
			v.updated_at,
			v.deleted_at,
			v.stock_qty,
			v.reserved_qty,

			COALESCE(
				attrs.attributes,
				'[]'::json
			) AS product_variant_attributes,

			COALESCE(
				files.files,
				'[]'::json
			) AS product_files

		FROM product_variants v

		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'id', a.id,
					'variant_id', a.variant_id,
					'name', a.name,
					'value', a.value
				)
				ORDER BY a.id
			) AS attributes
			FROM product_variant_attributes a
			WHERE a.variant_id = v.id
		) attrs ON TRUE

		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'file_id', pf.product_file_id,
					'product_id', pf.product_id,
					'variant_id', pf.variant_id,
					'content_type', pf.content_type,
					'position', pf.position,
					'is_primary', pf.is_primary
				)
				ORDER BY pf.position
			) AS files
			FROM product_file pf
			WHERE pf.product_id = v.product_id
			  AND pf.variant_id = v.id
		) files ON TRUE

		WHERE v.product_id = $1
		  AND v.deleted_at IS NULL
		ORDER BY v.id;
	`

	var rows []variantRow

	if err := r.roDb.SelectContext(ctx, &rows, q, productID); err != nil {
		return nil, logger.Errorf(
			ctx,
			err,
			"failed to fetch product variants for product %d",
			productID,
		)
	}

	variants := make(model.ProductVariants, 0, len(rows))

	for _, row := range rows {
		variant := row.ProductVariant

		if len(row.AttrsJSON) > 0 {
			if err := json.Unmarshal(
				row.AttrsJSON,
				&variant.ProductVariantAttributes,
			); err != nil {
				return nil, logger.Errorf(
					ctx,
					err,
					"failed to unmarshal attributes for variant %d",
					variant.ID,
				)
			}
		}

		if len(row.FilesJSON) > 0 {
			if err := json.Unmarshal(
				row.FilesJSON,
				&variant.ProductFiles,
			); err != nil {
				return nil, logger.Errorf(
					ctx,
					err,
					"failed to unmarshal files for variant %d",
					variant.ID,
				)
			}
		}

		variants = append(variants, variant)
	}

	return variants, nil
}

func (r *productVariantsRepo) GetVariantsByID(ctx context.Context, variantId uint) (*model.ProductVariant, error) {
	const q = `SELECT v.id, v.product_id, v.sku, v.price_amount, v.created_at, v.updated_at, v.deleted_at,
		v.stock_qty, v.reserved_qty,
		COALESCE(
			json_agg(
				json_build_object(
					'id', a.id,
					'variant_id', a.variant_id,
					'name', a.name,
					'value', a.value
				)
			) FILTER (WHERE a.id IS NOT NULL),
	'[]'::json
	) AS product_variant_attributes
	FROM product_variants v
	LEFT JOIN product_variant_attributes a ON a.variant_id = v.id
	WHERE v.id = $1
	AND v.deleted_at IS NULL
	GROUP BY
	v.id, v.product_id, v.sku, v.price_amount, v.created_at, v.updated_at, v.deleted_at,
		v.stock_qty, v.reserved_qty;`

	var row variantRow
	if err := r.roDb.GetContext(ctx, &row, q, variantId); err != nil {
		err = logger.Errorf(ctx, err, "failed to fetch product variant %d", variantId)
		return nil, err
	}
	v := row.ProductVariant
	// attributes
	if len(row.AttrsJSON) > 0 {
		_ = json.Unmarshal(row.AttrsJSON, &v.ProductVariantAttributes)
	}
	return &v, nil
}

func (r *productVariantsRepo) GetVariantsByIDs(ctx context.Context, variantIDs []uint) (model.ProductVariants, error) {
	var variants model.ProductVariants
	if len(variantIDs) == 0 {
		return variants, nil
	}

	q := `SELECT id, product_id, sku, price_amount, created_at, updated_at, deleted_at, stock_qty, reserved_qty
		FROM product_variants
		WHERE id = ANY($1)
		  AND deleted_at IS NULL;`

	if err := r.roDb.SelectContext(ctx, &variants, q, pq.Array(variantIDs)); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get product variants by ids: %v", variantIDs)
	}

	return variants, nil
}

func (r *productVariantsRepo) GetAttrByVariantsIDs(ctx context.Context, variantIDs []uint) (model.ProductVariantAttributes, error) {
	var attrs model.ProductVariantAttributes

	if len(variantIDs) == 0 {
		return attrs, nil
	}

	q := `SELECT * FROM product_variant_attributes WHERE variant_id = ANY($1);`

	if err := r.roDb.SelectContext(ctx, &attrs, q, pq.Array(variantIDs)); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get variant attributes")
	}
	return attrs, nil
}

func (r *productVariantsRepo) ReserveStockForCheckout(ctx context.Context, tx *sqlx.Tx, variantID uint, qty uint) error {
	const q = `UPDATE product_variants
        SET reserved_qty = reserved_qty + $2,
            updated_at = now()
        WHERE id = $1
          AND (stock_qty - reserved_qty) >= $2`

	res, err := tx.ExecContext(ctx, q, variantID, qty)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to reserve stock for variant %d", variantID)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return logger.Errorf(ctx, err, "failed to get rows affected for variant %d", variantID)
	}
	if rows == 0 {
		return logger.Errorf(ctx, nil, "insufficient stock for variant %d", variantID)
	}
	return nil
}
