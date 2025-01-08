package plog

import (
	"context"
	"time"
)

type Logger interface {
	Out(Output) Logger
	Type(any) Logger
	Scope(string, ...any) Logger
	Set(ctx *context.Context) Logger
	Msg(string) Logger
	P() Logger
	F() Logger
	E() Logger
	W() Logger
	I() Logger
	D() Logger
	Flush()
	Printf(string, ...any)
	Println(...any)
	Copy() Logger
}

type Output func(Log)

type Log struct {
	Level   Level
	Pid     int
	Scopes  []string
	Time    time.Time
	Message string
	Stack   []byte
	Errors  []error
}

type Level int

const (
	Panic Level = iota
	Fatal
	Error
	Warning
	Info
	Debug
	all
)
