package categorize

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"groceryspend.io/server/utils"
)

type CategoryExternalService interface {
	GetCategoryForItem(items []string, target *map[string]string) error
}

type DefaultCategoryExternalService struct {
	hostname string
	path     string
}

func NewDefaultCategoryExternalService() *DefaultCategoryExternalService {
	return &DefaultCategoryExternalService{
		hostname: utils.GetOsValue("CATEGORIZE_HOST"),
		path:     utils.GetOsValue("CATEGORIZE_PATH"),
	}
}

func (s *DefaultCategoryExternalService) GetCategoryForItem(items []string, target *map[string]string) error {

	itemsJson, _ := json.Marshal(items)
	body := strings.NewReader(string(itemsJson))
	// make HTTP call to category service
	resp, err := http.Post(fmt.Sprintf("%v/%v", s.hostname, s.path), "application/json", body)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
