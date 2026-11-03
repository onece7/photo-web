package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Phone        string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
