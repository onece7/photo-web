package main

import (
	"fmt"
	"photo-web/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// 自动迁移
	if err := db.AutoMigrate(
		&models.User{},
		&models.ImageStyle{},
		&models.GenerationRecord{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败：%w", err)
	}
	return nil
}
