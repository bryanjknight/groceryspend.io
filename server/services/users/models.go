package users

import (
	"github.com/google/uuid"
)

// User represents the canonical user
type User struct {
	UserUUID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	AuthProviderID string    `gorm:"notNull;index:auth_provider_idx"`
	AuthSpecificID string    `gorm:"notNull;index:auth_provider_idx"`
}
