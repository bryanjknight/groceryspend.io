package receipts

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

const testHTMLContent = `
<html>
<body>
<div id="1">
	<div id="2">
		<a id="3" href="https://example.com">Hi</a>
	</div>
</div>
</body>
</html>`

func TestElementById(t *testing.T) {
	testHTML, err := html.Parse(strings.NewReader(testHTMLContent))
	if err != nil {
		panic("Failed to parse html")
	}

	type test struct {
		input         string
		outputTagName string
	}

	tests := []test{
		{input: "1", outputTagName: "div"},
		{input: "2", outputTagName: "div"},
	}

	for _, tc := range tests {
		got := GetElementByID(testHTML, tc.input)
		if got == nil || (got != nil && got.Data != tc.outputTagName) {
			t.Errorf("expected: %v, got: %v", tc.outputTagName, got)
		}
	}
}

func TestElementByIdFailConditions(t *testing.T) {
	testHTML, err := html.Parse(strings.NewReader(testHTMLContent))
	if err != nil {
		panic("Failed to parse html")
	}

	type test struct {
		input string
	}

	tests := []test{
		{input: "a"},
		{input: "b"},
	}

	for _, tc := range tests {
		got := GetElementByID(testHTML, tc.input)
		if got != nil {
			t.Errorf("expected nil, got: %v", got)
		}
	}
}
func TestElementByTagName(t *testing.T) {
	testHTML, err := html.Parse(strings.NewReader(testHTMLContent))
	if err != nil {
		panic("Failed to parse html")
	}

	type test struct {
		input string
		ids   []string
	}

	tests := []test{
		{input: "div", ids: []string{"1", "2"}},
		{input: "a", ids: []string{"3"}},
		{input: "abc", ids: []string{}},
	}

	for _, tc := range tests {
		got := GetElementsByTagName(testHTML, tc.input)

		// get id attribute
		gotIds := []string{}
		for _, n := range got {
			id, _ := GetAttribute(n, "id")
			gotIds = append(gotIds, id)
		}

		if !reflect.DeepEqual(tc.ids, gotIds) {
			t.Errorf("expected: %v, got: %v", tc.ids, gotIds)
		}
	}
}
