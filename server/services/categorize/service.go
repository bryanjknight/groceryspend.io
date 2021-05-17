package categorize

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"groceryspend.io/server/utils"
)

// CategoryExternalService is a service that can provide categorizing of items
type CategoryExternalService interface {
	GetCategoryForItems(items []string, target *map[string]string) error
}

// DefaultCategoryExternalService is the default service
type DefaultCategoryExternalService struct {
	baseURL string
}

// NewDefaultCategoryExternalService creates a new http wrapper around the http endpoint
func NewDefaultCategoryExternalService() *DefaultCategoryExternalService {
	return &DefaultCategoryExternalService{
		baseURL: fmt.Sprintf("%s/%s", utils.GetOsValue("CATEGORIZE_HOST"), utils.GetOsValue("CATEGORIZE_PATH")),
	}
}

// GetCategoryForItems takes an array of items and returns a mapping of item to category. If duplicates are
// in the list, the last one wins
func (s *DefaultCategoryExternalService) GetCategoryForItems(items []string, target *map[string]string) error {

	itemsJSON, _ := json.Marshal(items)
	body := strings.NewReader(string(itemsJSON))
	// make HTTP call to category service
	resp, err := http.Post(s.baseURL, "application/json", body)

	if err != nil {
		println(fmt.Sprintf("Failed to get response from prediction service, %s", err.Error()))
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected response 200, got %s", resp.Status)
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
