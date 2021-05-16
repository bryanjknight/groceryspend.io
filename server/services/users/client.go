package users

// Client is an internal mechanism for accessing the user service
type Client interface {
	LookupUserByAuthProvider(authProvider string, authSpecificID string) (*User, error)
}

// DefaultClient is the default client that uses a user repo for access
type DefaultClient struct {
	repo UserRepository
}

// NewDefaultClient creates a default client. Currently defaults to using postgres user repo
func NewDefaultClient() *DefaultClient {

	// if we wanted to experiment with different database, here's how we would do the switch
	// for now, just assume the postgres repo
	pgRepo := NewPostgresUserRepository()

	return &DefaultClient{repo: pgRepo}
}

// LookupUserByAuthProvider finds (or creates) a user based on the auth provider and auth specific ID for the user
func (c *DefaultClient) LookupUserByAuthProvider(authProvider string, authSpecificID string) (*User, error) {
	user, err := c.repo.GetOrCreateUserByAuthProviderID(authProvider, authSpecificID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
