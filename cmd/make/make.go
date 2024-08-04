package main

import (
	"log"
	"strconv"
)

type Make int

const (
	None   Make = 0
	System Make = 1 << (iota - 1)
	Libs
	Apps
	Dev
	Portal
	Installer
	All = System | Libs | Apps | Dev | Portal | Installer
)

var mappings = map[rune]Make{
	'l': Libs,
	'a': Apps,
	'd': Dev,
	'p': Portal,
	'i': Installer,
}

type Install struct {
	root string
}

func NewInstall(root string) *Install {
	return &Install{root: root}
}

func ParseArgs(args []string) (jobs Make) {
	for _, arg := range args {
		if i, err := strconv.Atoi(arg); err == nil {
			jobs += Make(i)
			continue
		}
		for _, r := range []rune(arg) {
			jobs += mappings[r]
		}
	}
	if jobs == None {
		jobs = All
	}
	log.Println("parsed jobs", jobs)
	return
}

func (d *Install) Run(make Make, goos []string) {
	resolveVersion()
	defer clearVersion()
	if make&System == System {
		// no-op
	}
	if make&Libs == Libs {
		log.Println(" * js libs")
		d.buildJsLibs()
	}
	if make&Apps == Apps {
		log.Println(" * embed apps")
		d.buildEmbedApps(goos...)
	}
	if make&Dev == Dev {
		log.Println(" * portal dev")
		goPortalDev.Install()
		goPortalDevExec.Install()
		goPortalDevGo.Install()
		goPortalDevGoja.Install()
		goPortalDevWails.Install()
	}
	if make&Portal == Portal {
		log.Println(" * portal")
		goPortal.Install()
		goPortalApp.Install()
		goPortalAppGoja.Install()
		goPortalAppWails.Install()
	}
	if make&Installer == Installer {
		log.Println(" * installer")
		d.buildInstaller(goos...)
	}
}
