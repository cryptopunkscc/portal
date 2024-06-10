package target

import (
	"strings"
)

type Port struct {
	Base string
	Name string
}

func (p Port) Copy(base string) Port {
	p.Base = base
	return p
}

func (p Port) Route(name string) Port {
	p.Name = name
	return p
}

func (p Port) Target(portal Portal) Port {
	p.Base = portal.Manifest().Package
	return p
}

func (p Port) String() string {
	var chunks []string
	if p.Base != "" {
		chunks = append(chunks, p.Base)
	}
	if p.Name != "" {
		chunks = append(chunks, p.Name)
	}
	return strings.Join(chunks, ".")
}
