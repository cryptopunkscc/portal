package registry

import (
	"math"
	"slices"
)

type Node[V any] struct {
	next     map[byte]*Node[V]
	value    V
	empty    V
	dividers []byte
}

func New[V any](dividers ...byte) *Node[V] {
	return &Node[V]{
		next:     make(map[byte]*Node[V]),
		dividers: dividers,
	}
}

func (n *Node[V]) AddAll(callers map[string]V) *Node[V] {
	for key, value := range callers {
		n.Add(key, value)
	}
	return n
}

func (n *Node[V]) Add(path string, value V) *Node[V] {
	path = n.filter(path)
	if len(path) == 0 {
		if n.IsEmpty() {
			n.value = value
		}
		return n
	}
	var next *Node[V]
	next, ok := n.next[path[0]]
	if !ok {
		next = New[V](n.dividers...)
		n.next[path[0]] = next
	}
	return next.Add(path[1:], value)
}

func (n *Node[V]) Set(path string, node *Node[V]) *Node[V] {
	path = n.filter(path)
	node.dividers = n.dividers
	if len(path) == 0 {
		*n = *node
		return n
	}
	var next *Node[V]
	next, ok := n.next[path[0]]
	if !ok {
		next = New[V](n.dividers...)
		n.next[path[0]] = next
	}
	next.Set(path[1:], node)
	return n
}

func (n *Node[V]) Get() V {
	return n.value
}

func (n *Node[V]) Fold(path string) (*Node[V], string) { return n.fold(path, 0, 0) }

func (n *Node[V]) fold(path string, offset int, skip int) (*Node[V], string) {
	index := offset + skip
	if len(path) == index {
		return n, path[offset:]
	}
	if n.skip(path[index]) {
		return n.fold(path, offset, skip+1)
	}
	nn, ok := n.next[path[index]]
	if !ok {
		return n, path[offset:]
	}
	return nn.fold(path, index+1, 0)
}

func (n *Node[V]) filter(route string) (out string) {
	for _, b := range []byte(route) {
		if !n.skip(b) {
			out = out + string(b)
		}
	}
	return
}

func (n *Node[V]) skip(b byte) bool {
	return slices.Contains(n.dividers, b)
}

func (n *Node[V]) HasNext() bool {
	return len(n.next) > 0
}

func (n *Node[V]) IsEmpty() bool {
	var value any = n.value
	var empty any = n.empty
	return value == empty
}

func (n *Node[V]) All() map[string]*Node[V] {
	m := make(map[string]*Node[V])
	n.list(math.MaxInt, nil, m)
	return m
}

func (n *Node[V]) Children() map[string]*Node[V] {
	m := make(map[string]*Node[V])
	n.list(1, nil, m)
	return m
}

func (n *Node[V]) list(depth int, str []byte, m map[string]*Node[V]) {
	for b, r := range n.next {
		d := depth
		s := append(str, b)
		if !r.IsEmpty() {
			m[string(s)] = r
			d -= 1
		}
		if d == 0 {
			continue
		}
		if len(r.next) > 0 {
			r.list(d, s, m)
		}
	}
}
