package categorize

import "fmt"

// Client is a a wrapper class around a categorize service
type Client interface {
	GetCategoryForItems(items []string, target *map[string]*Category) error
	GetCategoryByID(id uint) (*Category, error)
	GetAllCategories() ([]*Category, error)
}

// DefaultClient is the default client
type DefaultClient struct {
	service  CategoryExternalService
	idToName map[uint]string
}

// NewDefaultClient creates a new default client
func NewDefaultClient() *DefaultClient {
	service := NewDefaultCategoryExternalService()

	return &DefaultClient{
		service: service,
	}
}

// GetCategoryForItems returns categories for a list of item names
func (c *DefaultClient) GetCategoryForItems(items []string, target *map[string]*Category) error {

	err := c.service.GetCategoryForItems(items, target)
	if err != nil {
		return err
	}
	return nil
}

// GetCategoryByID gets the category by ID
func (c *DefaultClient) GetCategoryByID(id uint) (*Category, error) {

	cats, err := c.service.GetAllCategories()

	if err != nil {
		return nil, err
	}

	for _, cat := range cats {
		if id == cat.ID {
			return cat, nil
		}
	}

	return nil, fmt.Errorf("Did not find a category for value %v", id)
}

// GetAllCategories returns an array of categories
func (c *DefaultClient) GetAllCategories() ([]*Category, error) {
	return c.service.GetAllCategories()
}
