package component

import (
	"fmt"
	"io/ioutil"
	"math"
	"path"
	"strings"

	"gopkg.in/tent.v1/source"
)

// Picture represents an image.
type Picture struct {
	ID   string
	Data []byte
}

// Order returns math.MaxFloat64, Pictures are shown last.
func (*Picture) Order() float64 { return math.MaxFloat64 }

func (p Picture) String() string {
	return fmt.Sprintf("Picture:%s Size:%v", p.ID, len(p.Data))
}

// picParser is the Parser for Picture
type picParser struct{}

// Match tells if it's a Picture from the name.
func (picParser) Match(name string) bool {
	_, file := path.Split(name)
	ext := strings.ToLower(path.Ext(file))
	for _, e := range []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"} {
		if e == ext {
			return true
		}
	}
	return false
}

// Parse populates the Picture with Item contents.
func (picParser) Parse(root *Category, item source.Item) error {
	dir, file := path.Split(item.Name())
	contents, err := item.Content()
	if err != nil {
		return err
	}
	defer contents.Close()

	data, err := ioutil.ReadAll(contents)
	if err != nil {
		return err
	}
	cat := root.Ensure(dir)
	cat.Components = append(cat.Components, &Picture{
		ID:   file,
		Data: data,
	})
	return nil
}
