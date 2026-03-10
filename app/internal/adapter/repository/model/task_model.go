package model

import (
	"time"
)

// TaskModel is the GORM model for tasks. Kept in adapter layer to preserve clean architecture.
type TaskModel struct {
	ID               string    `gorm:"type:varchar(36);primaryKey"`
	UserID           string    `gorm:"type:varchar(36);not null;index"`
	Title            string     `gorm:"size:500;not null"`
	Content          string     `gorm:"type:text;not null"`
	DueDate          time.Time  `gorm:"type:date;not null"`
	Done             bool       `gorm:"not null;default:false"`
	RequestTimestamp time.Time  `gorm:"type:timestamptz;not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TableName overrides the table name.
func (TaskModel) TableName() string {
	return "tasks"
}
