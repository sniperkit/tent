package component

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"gopkg.in/tent.v1/source"
	yaml "gopkg.in/yaml.v2"
)

// Category is a parent node in the tree.
type Category struct {
	ID         string            `yaml:"id"`
	Index      float64           `yaml:"index"`
	Meta       map[string]string `yaml:",inline"`
	Sub        []Category        `yaml:"sub"`
	Components []Component       `yaml:"components"`
}

// Order implements the Component interface.
func (c *Category) Order() float64 { return c.Index }

func (c *Category) String() string {
	return fmt.Sprintf("Category:%s Idx:%v Meta:%v", c.ID, c.Index, c.Meta)
}

func (c *Category) sort() {
	sort.Slice(c.Sub, func(i, j int) bool {
		return c.Sub[i].Index < c.Sub[j].Index
	})
	sort.Slice(c.Components, func(i, j int) bool {
		return c.Components[i].Order() < c.Components[j].Order()
	})
	for i := range c.Sub {
		c.Sub[i].sort()
	}
}

// Ensure returns a child Category, creating all needed nodes.
func (c *Category) Ensure(dir string) *Category {
item:
	for _, id := range strings.FieldsFunc(dir, func(r rune) bool { return r == '/' }) {
		for i := range c.Sub {
			if c.Sub[i].ID == id {
				c = &c.Sub[i]
				continue item
			}
		}
		c.Sub = append(c.Sub, Category{ID: id})
		c = &c.Sub[len(c.Sub)-1]
	}
	return c
}

// Component represents a leaf node.
type Component interface {
	Order() float64
}

// catParser is the Parser for Category
type catParser struct{}

// Match tells if it's a Category from the name.
func (catParser) Match(name string) bool {
	_, file := path.Split(name)
	return file == catName
}

// Parse populates the Category with Item contents.
func (catParser) Parse(root *Category, item source.Item) error {
	contents, err := item.Content()
	if err != nil {
		return err
	}
	defer contents.Close()
	dir, _ := path.Split(item.Name())
	c := Category{ID: path.Base(dir)}
	if err := yaml.NewDecoder(contents).Decode(&c); err != nil {
		return err
	}
	if dir := path.Dir(path.Clean(dir)); dir != "." {
		root = root.Ensure(dir)
	}
	root.Sub = append(root.Sub, c)
	return nil
}
