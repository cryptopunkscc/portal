package rpc

import "io"

type Registry[V any] struct {
	next  map[byte]*Registry[V]
	value V
	empty V
}

func NewRegistry[V any]() *Registry[V] {
	return &Registry[V]{next: make(map[byte]*Registry[V])}
}

func (n *Registry[V]) Add(str string, v V) {
	if len(str) == 0 {
		n.value = v
		return
	}
	var next *Registry[V]
	next, ok := n.next[str[0]]
	if !ok {
		next = NewRegistry[V]()
		n.next[str[0]] = next
	}
	next.Add(str[1:], v)
}

func (n *Registry[V]) Get() V {
	return n.value
}

func (n *Registry[V]) Unfold(str string) (*Registry[V], string) {
	if len(str) == 0 {
		return n, str
	}
	nn, ok := n.next[str[0]]
	if !ok {
		return n, str
	}
	return nn.Unfold(str[1:])
}

func (n *Registry[V]) Scan(scanner io.ByteScanner) (rr *Registry[V], err error) {
	var b byte
	b, err = scanner.ReadByte()
	nn, ok := n.next[b]
	if !ok {
		return n, scanner.UnreadByte()
	}
	return nn.Scan(scanner)
}

func (n *Registry[V]) HasNext() bool {
	return len(n.next) > 0
}

func (n *Registry[V]) IsEmpty() bool {
	var value any = n.value
	var empty any = n.empty
	return value == empty
}

func (n *Registry[V]) All() map[string]V {
	m := make(map[string]V)
	n.all(nil, m)
	return m
}

func (n *Registry[V]) all(str []byte, m map[string]V) {
	for b, r := range n.next {
		s := append(str, b)
		if !r.IsEmpty() {
			m[string(s)] = r.value
		}
		if len(r.next) > 0 {
			r.all(s, m)
		}
	}
}
