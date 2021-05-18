package categorize

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCategoryAPI(t *testing.T) {

	// test data
	mockOutput := `
	{
		"EUROPEAN CELLO CUKES": {
			"id":   17,
			"name": "Produce"
		},
		"FAGE WHOLE 1/2 KILO": {
			"id":   14,
			"name": "Dairy"
		},
		"GR/HOUSE RED PEPPER 0.55 @ 1 lb / 2.99": {
			"id":   17,
			"name": "Produce"
		},
		"GRN MTN SALSA MILD": {
			"id":   9,
			"name": "Snacks & Candy"
		}
	}`

	expected := map[string]*Category{
		"EUROPEAN CELLO CUKES": {
			ID:   17,
			Name: "Produce",
		},
		"FAGE WHOLE 1/2 KILO": {
			ID:   14,
			Name: "Dairy",
		},
		"GR/HOUSE RED PEPPER 0.55 @ 1 lb / 2.99": {
			ID:   17,
			Name: "Produce",
		},
		"GRN MTN SALSA MILD": {
			ID:   9,
			Name: "Snacks & Candy",
		},
	}

	actual := make(map[string]*Category)
	err := json.NewDecoder(strings.NewReader(mockOutput)).Decode(&actual)

	if err != nil {
		t.Fatalf("Invalid mock output: %s", err.Error())
	}

	// check key length
	if len(actual) != len(expected) {
		t.Errorf("Expected %v values, got %v", len(expected), len(actual))
	}

	// for each item, check the category
	for key, value := range actual {
		expectedValue := expected[key]
		if value.ID != expectedValue.ID || value.Name != expectedValue.Name {
			t.Errorf("For item %s, expected %s, got %s", key, expectedValue, value)
		}
	}

}
