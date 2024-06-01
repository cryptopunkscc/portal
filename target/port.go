package target

import (
	"strings"
)

type Port struct {
	Prefix  []string
	Host    string
	Command string
}

func (p Port) Copy(host string) Port {
	p.Host = host
	return p
}

func (p Port) Cmd(command string) Port {
	p.Command = command
	return p
}

func (p Port) Target(portal Portal) Port {
	p.Host = portal.Manifest().Package
	return p
}

func (p Port) String() string {
	chunks := p.Prefix
	if p.Host != "" {
		chunks = append(chunks, p.Host)
	}
	if p.Command != "" {
		chunks = append(chunks, p.Command)
	}
	return strings.Join(chunks, ".")
}
