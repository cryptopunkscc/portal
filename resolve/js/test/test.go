package test

import "embed"

//go:embed dist
var distFS embed.FS

//go:embed project
var projectFS embed.FS
