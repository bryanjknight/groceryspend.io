package users

import (
	"github.com/google/uuid"
)

// User is a canonical user used by other services that is tied to an auth provider
type User struct {
	ID             uuid.UUID `json:"id"`
	AuthProviderID string    `json:"authProviderID" db:"auth_provider_id"`
	AuthSpecificID string    `json:"authSpecificID" db:"auth_specific_id"`
}
