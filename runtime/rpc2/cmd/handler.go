package cmd

import (
	"strings"
)

type Root Handler
type Handlers []Handler
type Params []Param

type Handler struct {
	Func   any
	Name   string
	Desc   string
	Params []Param
	Sub    []Handler
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
