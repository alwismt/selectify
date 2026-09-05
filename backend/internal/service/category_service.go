package service

import (
	"context"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type categoryService struct {
	categoryRepo repo.CategoryRepo
}

type CategoryService interface {
	GetCategories(ctx context.Context) (model.Categories, error)
	AddProductCategoryWithTx(ctx context.Context, tx *sqlx.Tx, productID uint, categoryIDs []uint) error
}

func NewCategoryService(categoryRepo repo.CategoryRepo) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetCategories(ctx context.Context) (model.Categories, error) {
	categories, err := s.categoryRepo.GetCategories(ctx)
	if err != nil {
		return categories, logger.Error(ctx, err, "failed to get categories")
	}

	return s.buildCategoryTree(categories), nil
}

func (s *categoryService) buildCategoryTree(categories model.Categories) model.Categories {
	categoryMap := make(map[uint]*model.Category, len(categories))
	roots := make(model.Categories, 0)

	for _, category := range categories {
		categoryMap[category.CategoryID] = category
	}

	for _, category := range categories {
		if category.ParentID == nil {
			roots = append(roots, category)
			continue
		}

		if parent, ok := categoryMap[*category.ParentID]; ok {
			parent.Children = append(parent.Children, category)
		}
	}

	return roots
}

func (s *categoryService) AddProductCategoryWithTx(ctx context.Context, tx *sqlx.Tx, productID uint, categoryIDs []uint) error {
	return s.categoryRepo.AddProductCategoriesWithTx(ctx, tx, productID, categoryIDs)
}
