package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"log"
	"path"
	"reflect"
)

type Project any

func Find[T target.Source](files fs.FS, dir string, filter ...T) (in <-chan T) {
	var filters []target.Source
	for _, f := range filter {
		v := reflect.ValueOf(f)
		if reflect.TypeOf(f).Kind() != reflect.Pointer {
			filters = append(filters, f)
		} else {
			filters = append(filters, v.Elem().Interface().(target.Source))
		}
	}
	var f T
	if len(filters) == 0 {
		filters = append(filters, f)
	}
	out := make(chan T)
	in = out
	go func() {
		defer close(out)
		for project := range FindAll(files, dir, filters...) {
			p, ok := project.(T)
			if !ok {
				p = reflect.ValueOf(project).Elem().Interface().(T)
			}
			out <- p
		}
	}()
	return
}

func FindAll(files fs.FS, dir string, filter ...target.Source) (in <-chan target.Source) {
	out := make(chan target.Source)
	in = out
	if len(filter) == 0 {
		filter = append(filter, matchAll)
	}
	go func() {
		defer close(out)
		_ = fs.WalkDir(files, dir, func(src string, d fs.DirEntry, err error) error {
			if err != nil {
				return fs.SkipAll
			}
			if d.Name() == "node_modules" {
				return fs.SkipDir
			}

			if path.Ext(d.Name()) == ".portal" && d.Type().IsRegular() {
				log.Println("portal bundle detected: ", dir, files, src)

				module := NewModuleFS(src, files)
				bundle, err := module.Bundle()
				if err != nil {
					log.Println("bundle error:", err)
				}
				out <- bundle
				return nil
			}

			sub, err := fs.Sub(files, src)
			if err != nil {
				return err
			}
			module := NewModuleFS(src, sub)
			if nodeModule, err := module.NodeModule(); err == nil {

				if portalModule, err := nodeModule.PortalNodeModule(); err == nil {
					current := *portalModule
					for _, t := range filter {
						if isSameType(t, current) {
							log.Println("portal module detected: ", src)
							out <- &current
							return fs.SkipDir
						}
					}
				}
				current := *nodeModule
				for _, t := range filter {
					if isSameType(t, current) {
						log.Println("node module detected: ", src)
						out <- &current
						return nil
					}
				}
				return fs.SkipDir
			}
			if rawModule, err := module.PortalRawModule(); err == nil {
				current := *rawModule
				for _, t := range filter {
					if isSameType(t, current) {
						log.Println("raw module detected: ", src)
						out <- &current
						return fs.SkipDir
					}
				}
			}
			return nil
		})
	}()
	return
}

var matchAll = NewModuleFS("%all-matcher%", nil)

func isSameType(target target.Source, current target.Source) bool {
	if target == matchAll {
		return true
	}
	return reflect.TypeOf(current).AssignableTo(reflect.TypeOf(target))
}
