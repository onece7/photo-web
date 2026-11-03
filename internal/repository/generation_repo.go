package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"photo-web/internal/models"
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
	err := gorm.G[models.GenerationRecord](r.db).Create(ctx, record)
	return err

}

func (r *generationRepo) Update(ctx context.Context, record models.GenerationRecord) error {
	_, err := gorm.G[models.GenerationRecord](r.db).Where("id = ?", record.ID).Select("*").Updates(ctx, record)
	return err
}

func (r *generationRepo) FindByID(ctx context.Context, id uuid.UUID) (models.GenerationRecord, error) {
	record, err := gorm.G[models.GenerationRecord](r.db).Preload("Style", nil).Where("id = ?", id).First(ctx)
	if err != nil {
		return models.GenerationRecord{}, err
	}
	return record, nil
}

func (r *generationRepo) FindByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.GenerationRecord, int64, error) {
	query := gorm.G[models.GenerationRecord](r.db).Where("user_id = ?", userID)

	total, err := query.Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	records, err := query.Preload("Style", nil).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(ctx)

	return records, total, err
}

func (r *generationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[models.GenerationRecord](r.db).Where("id = ?", id).Delete(ctx)
	return err
}

func (r *generationRepo) CountPendingByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := gorm.G[models.GenerationRecord](r.db).
		Where("user_id = ? AND status IN ?", userID, []string{models.StatusPending, models.StatusProcessing}).
		Count(ctx, "*")
	return count, err
}
