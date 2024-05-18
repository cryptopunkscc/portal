package resolve

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"path"
)

func Source(module target.Source) (result target.Source, err error) {
	if path.Base(module.Path()) == "node_modules" {
		return nil, fs.SkipDir
	}
	module = module.Lift()
	bundle, err := project.ResolveBundle(module)
	if err == nil {
		result = bundle
		return
	}
	if path.Ext(module.Path()) != "" && module.Path() != "." {
		err = nil
		return
	}
	nodeModule, err := project.ResolveNodeModule(module)
	if err == nil {
		if result, err = project.ResolvePortalNodeModule(nodeModule); err == nil {
			return
		}
		result = nodeModule
		err = nil
		return
	}
	result, err = project.ResolvePortalRawModule(module)
	if err == nil {
		err = fs.SkipDir
		return
	}
	err = nil
	return
}
