package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
)

type PortalDev struct {
	root    string
	modules []project.NodeModule
}

func NewPortalDev(root string) *PortalDev {
	return &PortalDev{root: root}
}

type Dependency int

const (
	None   = Dependency(0x0)
	System = Dependency(0x1)
	Libs   = Dependency(0x2)
	Apps   = Dependency(0x4)
	Dev    = Dependency(0x8)
	All    = System | Libs | Apps | Dev
)

func (d *PortalDev) Install(deps ...Dependency) {
	dep := None
	for _, v := range deps {
		dep = dep | v
	}
	if dep == None {
		dep = All
	}
	log.Println("Portal dev installer")
	if dep&System == System {
		log.Println(" * native")
		checkGo()
		checkNpm()
		installAstral()
		installWails()
		installApt()
	}
	if dep&Libs == Libs {
		log.Println(" * js libs")
		d.buildJsLibs()
	}
	if dep&Apps == Apps {
		log.Println(" * js apps")
		d.BuildJsApps()
	}
	if dep&Dev == Dev {
		log.Println(" * portal dev")
		buildPortalDev()
	}
}
