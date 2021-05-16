package users

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	// load the postgres driver
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"

	// load source file driver
	_ "github.com/golang-migrate/migrate/source/file"

	"groceryspend.io/server/utils"
)

// ############################## //
// ##        WARNING           ## //
// ## Update this to match the ## //
// ## desired database version ## //
// ## for this git commit      ## //
// ############################## //

// DatabaseVersion is the desired database version for this git commit
const DatabaseVersion = 1

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

	// run migration
	migrationPath := "file://./services/users/db/migration"
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		log.Fatalf("Unable to get migration instance: %s", err)
	}

	err = m.Migrate(DatabaseVersion)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Database migration failed: %s", err)
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

	return &user, nil
}
