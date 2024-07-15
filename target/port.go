package target

import (
	"strings"
)

var PortPortal = NewPort("portal")
var PortOpen = PortPortal.Route("open")
var PortMsg = PortPortal.Route("broadcast")

type Port struct {
	Base string
	Name string
}

func NewPort(base string) Port {
	return Port{Base: base}
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
