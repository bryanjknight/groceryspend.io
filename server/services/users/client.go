package users

type Client interface {
	LookupUserByAuthProvider(authProvider string, authSpecificID string) (*User, error)
}

type DefaultClient struct {
	repo UserRepository
}

func NewDefaultClient() *DefaultClient {

	// if we wanted to experiment with different database, here's how we would do the switch
	// for now, just assume the postgres repo
	pgRepo := NewPostgresUserRepository()

	return &DefaultClient{repo: pgRepo}
}

func (c *DefaultClient) LookupUserByAuthProvider(authProvider string, authSpecificID string) (*User, error) {
	user, err := c.repo.GetOrCreateUserByAuthProviderID(authProvider, authSpecificID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
