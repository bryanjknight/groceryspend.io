package categorize

import "fmt"

// Category is a deparatment an item originates from (e.g. produce, deli, etc)
type Category struct {
	ID   uint   `json:"ID"`
	Name string `json:"Name"`
}

func (c *Category) String() string {
	return fmt.Sprintf("%s (%v)", c.Name, c.ID)
}
