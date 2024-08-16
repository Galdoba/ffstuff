package tag

import "fmt"

type collection struct {
	TagWithKey map[string]Tag
}

/*


 */

type Collection interface {
	AddTags(...Tag) error
	VerifyTags() error
	Base() string
	InputTag(string) string
	OutptTag(string) string
}

func NewCollection() *collection {
	return &collection{
		TagWithKey: make(map[string]Tag),
	}
}

func SearchTags(path string) (*collection, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *collection) AddTags(tags ...Tag) error {
	return fmt.Errorf("not implemented")
}

func (c *collection) VerifyTags() error {
	return fmt.Errorf("not implemented")
}

func (c *collection) Base() string {
	return "not implemented"
}

func (c *collection) InputTag(key string) string {
	return "not implemented"
}

func (c *collection) OutputTag(key string) string {
	return "not implemented"
}
