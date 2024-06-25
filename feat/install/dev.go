package install

import (
	"log"
)

type PortalDev struct {
	root string
}

func NewPortalDev(root string) *PortalDev {
	return &PortalDev{root: root}
}

type Dependency int

const (
	None   Dependency = 0
	System Dependency = 1 << (iota - 1)
	Libs
	Apps
	Dev
	Portal
	All = System | Libs | Apps | Dev | Portal
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
		d.buildJsApps()
	}
	if dep&Dev == Dev {
		log.Println(" * portal dev")
		buildPortalDev()
	}
	if dep&Portal == Portal {
		log.Println(" * portal")
		buildPortal()
	}
}
