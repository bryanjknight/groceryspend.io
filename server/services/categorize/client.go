package categorize

// Client is a a wrapper class around a categorize service
type Client interface {
	GetCategoryForItems(items []string, target *map[string]string) error
}

// DefaultClient is the default client
type DefaultClient struct {
	service CategoryExternalService
}

// NewDefaultClient creates a new default client
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		service: NewDefaultCategoryExternalService(),
	}
}

// GetCategoryForItems returns categories for a list of item names
func (c *DefaultClient) GetCategoryForItems(items []string, target *map[string]string) error {
	err := c.service.GetCategoryForItems(items, target)
	if err != nil {
		return err
	}
	return nil
}
