package users

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	// load the postgres river
	_ "github.com/lib/pq"

	"groceryspend.io/server/utils"
)

// UserRepository contains the common access patterns for canonical users
type UserRepository interface {
	GetOrCreateUserByAuthProviderID(authProvider string, authSpecificID string) (*User, error)
}

// PostgresUserRepository is an implementation of the user datastore using postgres
type PostgresUserRepository struct {
	db *sqlx.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepo
func NewPostgresUserRepository() *PostgresUserRepository {
	db, err := sqlx.Open("postgres", utils.GetOsValue("USERS_POSTGRES_CONN_STR"))
	if err != nil {
		panic(err)
	}
	retval := PostgresUserRepository{db: db}

	return &retval
}

// GetOrCreateUserByAuthProviderID look up user by auth provider id
func (r *PostgresUserRepository) GetOrCreateUserByAuthProviderID(authProvider string, authSpecificID string) (*User, error) {
	sql := `
	   INSERT INTO "user" (
	     auth_provider_id, auth_specific_id
	   ) VALUES (
	     $1, $2
	   )
	   ON CONFLICT(auth_provider_id, auth_specific_id)
	   -- TODO: perf impact of doing this
	   DO UPDATE SET auth_provider_id=EXCLUDED.auth_provider_id
	   RETURNING *;
	`
	// does a record exist?
	row := r.db.QueryRowxContext(context.Background(), sql, authProvider, authSpecificID)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}
	println(fmt.Sprintf("Extracted user: %s / %s / %s", user.ID.String(), user.AuthProviderID, user.AuthSpecificID))

	return &user, nil
}
