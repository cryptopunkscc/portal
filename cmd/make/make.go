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
		log.Println(" * native")
		checkGo()
		checkNpm()
		installAstral()
		installWails()
		installApt()
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
		installPortalDev()
	}
	if job&Portal == Portal {
		log.Println(" * portal")
		installPortal()
	}
	if job&Installer == Installer {
		log.Println(" * installer")
		d.buildInstaller()
	}
}
