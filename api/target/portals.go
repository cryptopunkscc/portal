package target

import (
	"slices"
)

func (p Portals[T]) SortBy(priority Priority) {
	if len(priority) > 1 {
		slices.SortFunc(p, func(a, b T) int {
			return priority.Get(a) - priority.Get(b)
		})
	}
}

func (p Portals[T]) Reduced() (reduced Portals[T]) {
	mem := make(map[string]T)
	for _, t := range p {
		if _, ok := mem[t.Manifest().Package]; ok {
			continue
		}
		mem[t.Manifest().Package] = t
		reduced = append(reduced, t)
	}
	return
}
