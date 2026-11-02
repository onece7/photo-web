package migration

import (
	"log"
	"log/slog"
	"photo-web/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	// 自动迁移
	if err := db.AutoMigrate(
		&models.User{},
		&models.ImageStyle{},
		&models.GenerationRecord{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
}

func SeedStyles(db *gorm.DB, logger *slog.Logger) {
	var count int64
	if err := db.Model(&models.ImageStyle{}).Count(&count).Error; err != nil {
		logger.Error("检查风格数据失败", slog.String("err", err.Error()))
		return
	}

	if count > 0 {
		return
	}

	styles := []models.ImageStyle{
		{
			Name: "anime", DisplayName: "动漫风格", Description: "日系动漫风格，色彩鲜艳，线条清晰",
			PresetParams: models.JSONMap{"negative_prompt": "realistic, photo, 3d render", "cfg_scale": 7.0, "steps": float64(20), "sampler": "Euler a"},
			SortOrder:    1, IsActive: true,
		},
		{
			Name: "realistic", DisplayName: "写实风格", Description: "照片级写实风格，细节丰富",
			PresetParams: models.JSONMap{"negative_prompt": "cartoon, anime, illustration", "cfg_scale": 7.5, "steps": float64(30), "sampler": "DPM++ 2M Karras"},
			SortOrder:    2, IsActive: true,
		},
		{
			Name: "watercolor", DisplayName: "水彩风格", Description: "柔和的水彩画风格",
			PresetParams: models.JSONMap{"negative_prompt": "photo, realistic, 3d", "cfg_scale": 6.5, "steps": float64(25), "sampler": "Euler a"},
			SortOrder:    3, IsActive: true,
		},
		{
			Name: "oil_painting", DisplayName: "油画风格", Description: "古典油画风格，笔触明显",
			PresetParams: models.JSONMap{"negative_prompt": "photo, digital art, anime", "cfg_scale": 7.0, "steps": float64(30), "sampler": "DPM++ 2M Karras"},
			SortOrder:    4, IsActive: true,
		},
		{
			Name: "pixel_art", DisplayName: "像素风格", Description: "复古像素艺术风格",
			PresetParams: models.JSONMap{"negative_prompt": "realistic, smooth, high resolution", "cfg_scale": 7.0, "steps": float64(20), "sampler": "Euler a"},
			SortOrder:    5, IsActive: true,
		},
		{
			Name: "sketch", DisplayName: "素描风格", Description: "黑白素描/线稿风格",
			PresetParams: models.JSONMap{"negative_prompt": "color, painting, realistic", "cfg_scale": 6.0, "steps": float64(20), "sampler": "Euler a"},
			SortOrder:    6, IsActive: true,
		},
	}
	if err := db.Create(&styles).Error; err != nil {
		logger.Error("插入风格种子失败", slog.String("err", err.Error()))
		return
	}
	logger.Info("已插入6条预设风格种子")
}
