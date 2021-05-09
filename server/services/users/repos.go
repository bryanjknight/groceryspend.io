package users

import (
	"context"
	"database/sql"

	// load the postgres river
	_ "github.com/lib/pq"

	"groceryspend.io/server/utils"
)

// UserRepository contains the common access patterns for canonical users
type UserRepository interface {
	GetOrCreateUserByAuthProviderID(authProvider string, authSpecificID string) (*User, error)
	// DisableUser(userUUID string) error
}

// PostgresUserRepository is an implementation of the user datastore using postgres
type PostgresUserRepository struct {
	q *Queries
}

// NewPostgresUserRepository creates a new PostgresUserRepo
func NewPostgresUserRepository() *PostgresUserRepository {
	db, err := sql.Open("postgres", utils.GetOsValue("USERS_POSTGRES_CONN_STR"))
	if err != nil {
		panic(err)
	}
	retval := PostgresUserRepository{&Queries{db: db}}

	return &retval
}

// GetOrCreateUserByAuthProviderID look up user by auth provider id
func (r *PostgresUserRepository) GetOrCreateUserByAuthProviderID(authProvider string, authSpecificID string) (*User, error) {

	// does a record exist?
	user, err := r.q.GetorCreateUserByAuthProviderId(context.Background(), GetorCreateUserByAuthProviderIdParams{
		AuthProviderID: authProvider,
		AuthSpecificID: authSpecificID,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}
