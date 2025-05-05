package cmd

import (
	"strings"
)

type Root Handler
type Handlers []Handler
type Params []Param

type Handler struct {
	Func   any       `json:"func,omitempty" yaml:"func,omitempty"`
	Name   string    `json:"name,omitempty" yaml:"name,omitempty"`
	Desc   string    `json:"description,omitempty" yaml:"description,omitempty"`
	Params []Param   `json:"params,omitempty" yaml:"params,omitempty"`
	Sub    []Handler `json:"sub,omitempty" yaml:"sub,omitempty"`
}

type Param struct {
	Name string
	Type string
	Desc string
}

func (h *Handler) AddSub(handlers ...Handler) {
	h.Sub = append(h.Sub, handlers...)
}

func (h *Handler) Names() []string {
	return strings.Split(h.Name, " ")
}

func (h Handlers) Plus(handlers ...Handler) Handlers {
	return append(h, handlers...)
}
