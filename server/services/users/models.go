package users

import (
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// User represents the canonical user
type User struct {
	UserUUID       string `gorm:"primaryKey"`
	AuthProviderID string `gorm:"notNull"`
	AuthSpecificID string `gorm:"notNull"`
}

// BeforeCreate before creating model, set the UserUUID to a generated UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	tx.Model(u).UpdateColumn("UserUUID", uuid.NewString())
	return nil
}
