package repository

import (
	"context"
	"photo-web/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StyleRepository 风格数据访问接口
type StyleRepository interface {
	FindAll(ctx context.Context) ([]models.ImageStyle, error)
	FindActive(ctx context.Context) ([]models.ImageStyle, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.ImageStyle, error)
	FindByName(ctx context.Context, name string) (*models.ImageStyle, error)
}

type styleRepo struct {
	db *gorm.DB
}

// NewStyleRepository 创建风格 Repository
func NewStyleRepository(db *gorm.DB) StyleRepository {
	return &styleRepo{db: db}
}

func (r *styleRepo) FindAll(ctx context.Context) ([]models.ImageStyle, error) {
	var styles []models.ImageStyle
	err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&styles).Error
	return styles, err
}

func (r *styleRepo) FindActive(ctx context.Context) ([]models.ImageStyle, error) {
	var styles []models.ImageStyle
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&styles).Error
	return styles, err
}

func (r *styleRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.ImageStyle, error) {
	var style models.ImageStyle
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&style).Error
	if err != nil {
		return nil, err
	}
	return &style, nil
}

func (r *styleRepo) FindByName(ctx context.Context, name string) (*models.ImageStyle, error) {
	var style models.ImageStyle
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&style).Error
	if err != nil {
		return nil, err
	}
	return &style, nil
}
