package rpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry[any]()
	r.Add("aaa", "a")
	r.Add("bbb", "b")

	s, v := r.Unfold("aaa")
	assert.Equal(t, "a", s.Get())
	assert.Equal(t, "", v)

	s, v = r.Unfold("aaabb")
	assert.Equal(t, "a", s.Get())
	assert.Equal(t, "bb", v)

	s, v = r.Unfold("bbb")
	assert.Equal(t, "b", s.Get())
	assert.Equal(t, "", v)

	s, v = r.Unfold("ccc")
	assert.Equal(t, nil, s.Get())
	assert.Equal(t, "ccc", v)
}

func TestRegistry_All(t *testing.T) {
	r := NewRegistry[any]()
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
