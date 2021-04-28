package categorize

type Client interface {
	GetCategoryForItems(items []string, target *map[string]string) error
}

type DefaultClient struct {
	service CategoryExternalService
}

func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		service: NewDefaultCategoryExternalService(),
	}
}

func (c *DefaultClient) GetCategoryForItems(items []string, target *map[string]string) error {
	err := c.service.GetCategoryForItems(items, target)
	if err != nil {
		return err
	}
	return nil
}
