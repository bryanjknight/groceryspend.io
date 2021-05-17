package categorize

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"groceryspend.io/server/utils"
)

type CategoryExternalService interface {
	GetCategoryForItems(items []string, target *map[string]string) error
}

type DefaultCategoryExternalService struct {
	baseUrl string
}

func NewDefaultCategoryExternalService() *DefaultCategoryExternalService {
	return &DefaultCategoryExternalService{
		baseUrl: fmt.Sprintf("%s/%s", utils.GetOsValue("CATEGORIZE_HOST"), utils.GetOsValue("CATEGORIZE_PATH")),
	}
}

func (s *DefaultCategoryExternalService) GetCategoryForItems(items []string, target *map[string]string) error {

	itemsJson, _ := json.Marshal(items)
	body := strings.NewReader(string(itemsJson))
	// make HTTP call to category service
	// NOTE we need to check the response code
	println(s.baseUrl)
	resp, err := http.Post(s.baseUrl, "application/json", body)

	if err != nil {
		println(fmt.Sprintf("Failed to get response from prediction service, %s", err.Error()))
		return err
	}

	println(fmt.Sprintf("Response Code: %s", resp.Status))

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
