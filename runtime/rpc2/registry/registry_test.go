package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistry(t *testing.T) {
	r := New[any]()
	r.Add("aaa", "a")
	r.Add("bbb", "b")

	s, v := r.Fold("aaa")
	assert.Equal(t, "a", s.Get())
	assert.Equal(t, "", v)

	s, v = r.Fold("aaabb")
	assert.Equal(t, "a", s.Get())
	assert.Equal(t, "bb", v)

	s, v = r.Fold("bbb")
	assert.Equal(t, "b", s.Get())
	assert.Equal(t, "", v)

	s, v = r.Fold("ccc")
	assert.Equal(t, nil, s.Get())
	assert.Equal(t, "ccc", v)
}

func TestRegistry_All(t *testing.T) {
	r := New[any]()
	r.Add("a", "a")
	r.Add("aaa", "aaa")
	r.Add("bbb", "bbb")

	expected := map[string]any{
		"a":   "a",
		"aaa": "aaa",
		"bbb": "bbb",
	}
	actual := r.All()
	assert.Equal(t, expected, actual)
}

func TestRegistry_Children(t *testing.T) {
	r := New[any]()
	r.Add("a", "a")
	r.Add("aaa", "aaa")
	r.Add("bbb", "bbb")

	expected := map[string]any{
		"a":   "a",
		"bbb": "bbb",
	}
	actual := r.Children()
	assert.Equal(t, expected, actual)
}

func TestRegistry_should_populate_dividers(t *testing.T) {
	r := New[any]('a')
	r.Add("foo", struct{}{})
	r.dividers[0] = 'b'
	rr, _ := r.Fold("foo")
	assert.ElementsMatch(t, r.dividers, rr.dividers)
}

func TestRegistry_should_not_populate_dividers(t *testing.T) {
	r := New[any]('.', ' ')
	r.Add("foo", 0)
	r.Add("foo.bar", 1)
	rr, args := r.Fold("foo.bar asd")
	assert.Equal(t, 1, rr.Get())
	assert.Equal(t, "asd", args)
}
