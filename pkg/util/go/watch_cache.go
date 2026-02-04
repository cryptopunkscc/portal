package golang

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type WatchCache struct {
	projectRoot string
	moduleName  string
	dirs        map[string]int
	files       map[string][]string
}

func NewWatchCache(root, moduleName string) *WatchCache {
	return &WatchCache{
		projectRoot: root,
		moduleName:  moduleName,
		dirs:        make(map[string]int),
		files:       make(map[string][]string),
	}
}

func (c *WatchCache) UpdateFile(abs string) (remove, add map[string]any) {
	oldImports, ok := c.files[abs]
	if ok {
		return
	}
	imports, err := Imports(abs)
	if err != nil {
		return
	}
	imports = filterByModule(imports, c.moduleName)
	imports = withProjectPath(imports, c.projectRoot)
	if !sameImports(imports, oldImports) {
		return
	}

	remove = c.RemoveFile(abs)

	c.files[abs] = imports
	add = map[string]any{}
	for _, s := range imports {
		for s2, a := range c.AddDir(s) {
			add[s2] = a
		}
	}
	return
}

func (c *WatchCache) AddFile(abs string) (dirs map[string]any) {
	if _, ok := c.files[abs]; ok {
		return
	}
	imports, err := Imports(abs)
	if err != nil {
		return
	}
	imports = filterByModule(imports, c.moduleName)
	imports = withProjectPath(imports, c.projectRoot)
	c.files[abs] = imports
	dirs = map[string]any{}
	for _, s := range imports {
		for s2, a := range c.AddDir(s) {
			dirs[s2] = a
		}
	}
	return
}

func (c *WatchCache) AddDir(dir string) (dirs map[string]any) {
	count := c.dirs[dir]
	c.dirs[dir] = count + 1
	if count > 0 {
		return make(map[string]any)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}
	dirs = map[string]any{}
	dirs[dir] = dir
	for _, entire := range entries {
		if entire.IsDir() {
			continue
		}
		if !strings.HasSuffix(entire.Name(), ".go") {
			continue
		}
		e := filepath.Join(dir, entire.Name())
		for s, a := range c.AddFile(e) {
			dirs[s] = a
		}
	}
	return
}

func (c *WatchCache) RemoveFile(file string) (dirs map[string]any) {
	dirs = make(map[string]any)
	imports, ok := c.files[file]
	if !ok {
		return
	}
	delete(c.files, file)
	for _, dir := range imports {
		for s, a := range c.PopDir(dir) {
			dirs[s] = a
		}
	}
	return
}

func (c *WatchCache) PopDir(dir string) (dirs map[string]any) {
	i := c.dirs[dir]
	if i == 0 {
		return
	}
	i = i - 1
	c.dirs[dir] = i
	gone := i == 0
	if !gone {
		return
	}
	return c.RemoveDir(dir)
}

func (c *WatchCache) RemoveDir(dir string) (dirs map[string]any) {
	dirs = make(map[string]any)
	dirs[dir] = dir
	delete(c.dirs, dir)
	for file, imports := range c.files {
		if strings.HasPrefix(file, dir) {
			delete(c.files, file)
			for _, dir := range imports {
				for s, a := range c.PopDir(dir) {
					dirs[s] = a
				}
			}
		}
	}
	return
}
