package users

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"groceryspend.io/server/utils"
)

// UserRepository contains the common access patterns for canonical users
type UserRepository interface {
	GetOrCreateUserByAuthProviderID(authProvider string, authSpecificID string) (*User, error)
	// DisableUser(userUUID string) error
}

// PostgresUserRepository is an implementation of the user datastore using postgres
type PostgresUserRepository struct {
	DbConnection *gorm.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepo
func NewPostgresUserRepository() *PostgresUserRepository {
	dbConn, err := gorm.Open(postgres.Open(utils.GetOsValue("USERS_POSTGRES_CONN_STR")), &gorm.Config{})

	if err != nil {
		panic("failed to connect to postgres db for users")
	}

	retval := PostgresUserRepository{DbConnection: dbConn}

	// TODO: this should be a script that runs as a different user. That way, the user running queries only
	//       has read/write but not create/delete permissions
	dbConn.AutoMigrate(&User{})

	return &retval
}

// GetOrCreateUserByAuthProviderID look up user by auth provider id
func (r *PostgresUserRepository) GetOrCreateUserByAuthProviderID(authProvider string, authSpecificID string) (*User, error) {

	var user User
	r.DbConnection.FirstOrCreate(&user, User{
		AuthProviderID: authProvider,
		AuthSpecificID: authSpecificID,
	})

	return &user, nil
}
