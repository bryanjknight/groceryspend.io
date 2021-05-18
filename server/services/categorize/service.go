package categorize

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kofalt/go-memoize"
	"groceryspend.io/server/utils"
)

// CategoryExternalService is a service that can provide categorizing of items
type CategoryExternalService interface {
	GetCategoryForItems(items []string, target *map[string]*Category) error
	GetAllCategories() ([]*Category, error)
}

// DefaultCategoryExternalService is the default service
type DefaultCategoryExternalService struct {
	baseURL        string
	catCache       memoize.Memoizer
	catIDTocatName map[int]string
}

// NewDefaultCategoryExternalService creates a new http wrapper around the http endpoint
func NewDefaultCategoryExternalService() *DefaultCategoryExternalService {

	catCache := memoize.NewMemoizer(90*time.Second, 10*time.Minute)

	return &DefaultCategoryExternalService{
		baseURL:  utils.GetOsValue("CATEGORIZE_HOST"),
		catCache: *catCache,
	}
}

// GetCategoryForItems takes an array of items and returns a mapping of item to category. If duplicates are
// in the list, the last one wins
func (s *DefaultCategoryExternalService) GetCategoryForItems(items []string, target *map[string]*Category) error {

	itemsJSON, _ := json.Marshal(items)
	body := strings.NewReader(string(itemsJSON))
	// make HTTP call to category service
	resp, err := http.Post(fmt.Sprintf("%s/%s", s.baseURL, "categorize"), "application/json", body)

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

// GetAllCategories fetch all categories
func (s *DefaultCategoryExternalService) GetAllCategories() ([]*Category, error) {

	baseURL := s.baseURL

	fetchAllCategoriesClosure := func() (interface{}, error) {
		return fetchAllCategories(baseURL)
	}

	categories, err, _ := s.catCache.Memoize("ALL_CATEGORIES", fetchAllCategoriesClosure)
	if err != nil {
		return nil, err
	}

	return categories.([]*Category), nil

}

func fetchAllCategories(baseURL string) ([]*Category, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, "categories"))

	if err != nil {
		println(fmt.Sprintf("Failed to get response from prediction service, %s", err.Error()))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected response 200, got %s", resp.Status)
	}

	defer resp.Body.Close()

	retval := []*Category{}
	json.NewDecoder(resp.Body).Decode(&retval)
	return retval, nil
}
