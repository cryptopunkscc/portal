package test

import _ "embed"

//go:embed test_build.yml
var BuildYml []byte

//go:embed test_dist.yaml
var DistYml []byte

//go:embed test_dev.yml
var DevYml []byte
