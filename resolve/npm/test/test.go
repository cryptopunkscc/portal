package test

import "embed"

//go:embed module
var moduleFS embed.FS

//go:embed project
var projectFS embed.FS
