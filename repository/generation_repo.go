package repository

import (
	"context"
	"photo-web/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenerationRepository 生成记录数据访问接口
type GenerationRepository interface {
	Create(ctx context.Context, record *models.GenerationRecord) error
	Update(ctx context.Context, record models.GenerationRecord) error
	FindByID(ctx context.Context, id uuid.UUID) (models.GenerationRecord, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.GenerationRecord, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CountPendingByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

type generationRepo struct {
	db *gorm.DB
}

// NewGenerationRepository 创建生成记录 Repository
func NewGenerationRepository(db *gorm.DB) GenerationRepository {
	return &generationRepo{db: db}
}

func (r *generationRepo) Create(ctx context.Context, record *models.GenerationRecord) error {
	return r.db.WithContext(ctx).Create(record).Error

}

func (r *generationRepo) Update(ctx context.Context, record models.GenerationRecord) error {
	return r.db.WithContext(ctx).
		Where("id = ?", record.ID).
		Select("*").
		Updates(&record).Error
}

func (r *generationRepo) FindByID(ctx context.Context, id uuid.UUID) (models.GenerationRecord, error) {
	var record models.GenerationRecord
	err := r.db.WithContext(ctx).Preload("Style").Where("id = ?", id).First(&record).Error
	if err != nil {
		return models.GenerationRecord{}, err
	}
	return record, nil
}

func (r *generationRepo) FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.GenerationRecord, int64, error) {
	var records []models.GenerationRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&models.GenerationRecord{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Style").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

func (r *generationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.GenerationRecord{}).Error
}

func (r *generationRepo) CountPendingByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.GenerationRecord{}).
		Where("user_id = ? AND status IN ?", userID, []string{models.StatusPending, models.StatusProcessing}).
		Count(&count).Error
	return count, err
}
