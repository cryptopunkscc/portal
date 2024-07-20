package main

import (
	"log"
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

type Install struct {
	root string
}

func NewInstall(root string) *Install {
	return &Install{root: root}
}

func (d *Install) Run(jobs ...Make) {
	job := None
	for _, v := range jobs {
		job = job | v
	}
	if job == None {
		job = All
	}
	resolveVersion()
	defer clearVersion()
	if job&System == System {
		// no-op
	}
	if job&Libs == Libs {
		log.Println(" * js libs")
		d.buildJsLibs()
	}
	if job&Apps == Apps {
		log.Println(" * js apps")
		d.buildJsApps()
	}
	if job&Dev == Dev {
		log.Println(" * portal dev")
		goPortalDev.Install()
		goPortalDevExec.Install()
		goPortalDevGo.Install()
		goPortalDevGoja.Install()
		goPortalDevWails.Install()
	}
	if job&Portal == Portal {
		log.Println(" * portal")
		goPortal.Install()
		goPortalApp.Install()
		goPortalAppGoja.Install()
		goPortalAppWails.Install()
		goPortalTray.Install()
	}
	if job&Installer == Installer {
		log.Println(" * installer")
		d.buildInstaller()
	}
}
