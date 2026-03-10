package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/task-manager/api/internal/adapter/repository/model"
)

// NewGormDB creates a GORM database connection and runs AutoMigrate.
func NewGormDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	if err := db.AutoMigrate(&model.UserModel{}, &model.TaskModel{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	return db, nil
}
