package model

// UserModel is the GORM model for users. Kept in adapter layer to preserve clean architecture.
type UserModel struct {
	ID           string `gorm:"type:varchar(36);primaryKey"`
	Email        string `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash string `gorm:"type:text;not null"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
}

// TableName overrides the table name.
func (UserModel) TableName() string {
	return "users"
}
