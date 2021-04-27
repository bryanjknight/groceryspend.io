package categorize

type Client interface {
	GetCategoryForItem(items []string, target *map[string]string) error
}

type DefaultClient struct {
	service CategoryExternalService
}

func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		service: NewDefaultCategoryExternalService(),
	}
}

func (c *DefaultClient) GetCategoryForItem(items []string, target *map[string]string) error {
	err := c.service.GetCategoryForItem(items, target)
	if err != nil {
		return err
	}
	return nil
}
