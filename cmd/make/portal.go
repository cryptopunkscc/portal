package main

import (
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
)

func installPortal() {
	if err := exec.Run(
		".", "go", "install",
		"-tags", "desktop,wv2runtime.download,production",
		"-ldflags", "-w -s",
		"./cmd/portal",
	); err != nil {
		log.Fatalln("portal dev install failed: ", err)
	}
	log.Println()
	log.Println("portal installed successfully")
}

func installPortalDev() {
	if err := exec.Run(
		".", "go", "install",
		"-tags", "dev",
		"./cmd/portal-dev",
	); err != nil {
		log.Fatalln("portal dev install failed: ", err)
	}
	log.Println()
	log.Println("portal dev installed successfully")
}

func buildPortal() {
	if err := exec.Run(
		".", "go", "build",
		"-tags", "desktop,wv2runtime.download,production,webkit2_41",
		"-ldflags", "-w -s",
		"-o", "./cmd/portal-installer/bin/",
		"./cmd/portal",
	); err != nil {
		log.Fatalln("portal dev install failed: ", err)
	}
	log.Println("Portal build succeed")
}

func buildPortalDev() {
	if err := exec.Run(
		".", "go", "build",
		"-tags", "dev,webkit2_41",
		"-o", "./cmd/portal-installer/bin/",
		"./cmd/portal-dev",
	); err != nil {
		log.Fatalln("portal dev build failed: ", err)
	}
	log.Println("Portal dev build succeed")
}

func buildPortalInstaller() {
	if err := exec.Run(
		".", "go", "build",
		"-o", "./bin/",
		"./cmd/portal-installer",
	); err != nil {
		log.Fatalln("portal installer build failed: ", err)
	}
	log.Println("Portal installer build succeed")
}
