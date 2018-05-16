package component

import (
	"io"
	"testing"

	"gopkg.in/tent.v1/source"
	yaml "gopkg.in/yaml.v2"
)

func TestDecodeNestedCategories(t *testing.T) {
	items := []source.MockItem{
		{ID: "a/.category.yml", Contents: "index: 2\nm: x"},
		{ID: "a/b/.category.yml", Contents: "index: 7\nm: y"},
		{ID: "a/b/d/.category.yml", Contents: "index: 20\nm: w"},
		{ID: "a/b/c/.category.yml", Contents: "index: 12\nm: z"},
	}
	r, err := Decode(&source.MockSource{Items: items})
	if err != nil {
		t.Fatal(err)
	}
	if l := len(r.Sub); l != 1 {
		t.Fatalf("Expected %d category, got %d", 1, l)
	}
	if id := r.Sub[0].ID; id != "a" {
		t.Fatalf("Expected %q category, got %q", "a", id)
	}
	if l := len(r.Sub[0].Sub); l != 1 {
		t.Fatalf("Expected %d category, got %d", 1, l)
	}
	if id := r.Sub[0].Sub[0].ID; id != "b" {
		t.Fatalf("Expected %q category, got %q", "b", id)
	}
	if l := len(r.Sub[0].Sub[0].Sub); l != 2 {
		t.Fatalf("Expected %d category, got %d", 2, l)
	}
	if id := r.Sub[0].Sub[0].Sub[0].ID; id != "c" {
		t.Fatalf("Expected %q category, got %q", "c", id)
	}
	if id := r.Sub[0].Sub[0].Sub[1].ID; id != "d" {
		t.Fatalf("Expected %q category, got %q", "d", id)
	}
}

func TestDecodeComponent(t *testing.T) {
	items := []source.MockItem{
		{ID: "m_hello.mock", Contents: `index: 10`},
	}
	r, err := Decode(&source.MockSource{Items: items}, mockDecoder{})
	if err != nil {
		t.Fatal(err)
	}

	if l := len(r.Components); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	m, ok := r.Components[0].(mockCmp)
	if !ok {
		t.Fatalf("Expected component to be a %T, got %T", new(mockCmp), r.Components[0])
	}
	if id := m.ID; id != "hello" {
		t.Fatalf("Expected %q segments, got %q", "a", id)
	}
	if i := m.Index; i != 10 {
		t.Fatalf("Expected %v index, got %v", 10, i)
	}
}

func TestDecodeNestedComponent(t *testing.T) {
	items := []source.MockItem{
		{ID: "cat/m_hello.mock", Contents: `index: 10`},
	}
	r, err := Decode(&source.MockSource{Items: items}, mockDecoder{})
	if err != nil {
		t.Fatal(err)
	}

	if l := len(r.Sub); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	if l := len(r.Sub[0].Components); l != 1 {
		t.Fatalf("Expected %d components, got %d", 1, l)
	}
	m, ok := r.Sub[0].Components[0].(mockCmp)
	if !ok {
		t.Fatalf("Expected component to be a %T, got %T", mockCmp{}, r.Sub[0].Components[0])
	}
	if id := m.ID; id != "hello" {
		t.Fatalf("Expected %q segments, got %q", "a", id)
	}
	if i := m.Index; i != 10 {
		t.Fatalf("Expected %v index, got %v", 10, i)
	}
}

type mockCmp struct {
	ID    string
	Index float64
}

func (m mockCmp) Order() float64 { return m.Index }

func (mockCmp) Encode() (io.Reader, error) { return nil, nil }

type mockDecoder struct{}

func (mockDecoder) Format() (string, []string) { return "m_", []string{".mock"} }

func (mockDecoder) Decode(id string, r io.Reader) (Component, error) {
	m := mockCmp{ID: id}
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
