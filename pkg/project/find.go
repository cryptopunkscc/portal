package project

import (
	"io/fs"
	"log"
	"reflect"
)

type Project any

func Find[P Project](files fs.FS, dir string) (in <-chan P) {
	var p P
	out := make(chan P)
	in = out
	go func() {
		defer close(out)
		for project := range FindAll(files, dir, p) {
			out <- project.(P)
		}
	}()
	return
}

func FindAll(files fs.FS, dir string, filter ...Project) (in <-chan Project) {
	out := make(chan Project)
	in = out
	if len(filter) == 0 {
		filter = append(filter, matchAll)
	}
	go func() {
		_ = fs.WalkDir(files, dir, func(src string, d fs.DirEntry, err error) error {
			if d.Name() == "node_modules" {
				return fs.SkipDir
			}
			sub, err := fs.Sub(files, src)
			if err != nil {
				return err
			}
			directory := Module{dir: src, files: sub}
			if nodeModule, err := directory.NodeModule(); err == nil {
				if portalModule, err := nodeModule.PortalNodeModule(); err == nil {
					current := *portalModule
					for _, target := range filter {
						if isSameType(target, current) {
							log.Println("portal module detected: ", src)
							out <- current
							return fs.SkipDir
						}
					}
				}
				current := *nodeModule
				for _, target := range filter {
					if isSameType(target, current) {
						log.Println("node module detected: ", src)
						out <- current
						return nil
					}
				}
				return fs.SkipDir
			}
			if rawModule, err := directory.PortalRawModule(); err == nil {
				current := *rawModule
				for _, target := range filter {
					if isSameType(target, current) {
						log.Println("raw module detected: ", src)
						out <- current
						return fs.SkipDir
					}
				}
			}
			return nil
		})
		close(out)
	}()
	return
}

var matchAll = struct{}{}

func isSameType(target any, current any) bool {
	if target == matchAll {
		return true
	}
	return reflect.TypeOf(current) == reflect.TypeOf(target)
}
