package golang

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type ImportRefs struct {
	Import string
	Refs   []string
}

func (i ImportRefs) String() string {
	return i.Import
}

func ListImports(abs string) (l []ImportRefs, err error) {
	i := imports{}
	if err = i.Collect(abs); err != nil {
		return
	}
	l = i.List()
	return
}

type imports struct {
	Project
	c   chan []string
	wg  sync.WaitGroup
	set map[string][]string
}

func (imp *imports) List() (l []ImportRefs) {
	for s, n := range imp.set {
		l = append(l, ImportRefs{Import: s, Refs: n})
	}
	//sort.Slice(l, func(i, j int) bool { return strings.Compare(l[i].Import, l[j].Import) > 1 })
	sort.Slice(l, func(i, j int) bool { return len(l[i].Refs) > len(l[j].Refs) })
	return
}

func (imp *imports) Collect(path string) (err error) {
	if err = imp.Resolve(path); err != nil {
		return
	}
	imp.c = make(chan []string)
	imp.set = make(map[string][]string)

	if filepath.IsAbs(path) {
		path = strings.TrimPrefix(path, imp.Dir)
		path = strings.TrimPrefix(path, "/")
	}
	path = filepath.Clean(path)
	go imp.collect(path)

	go func() {
		time.Sleep(10 * time.Millisecond)
		imp.wg.Wait()
		time.Sleep(100 * time.Millisecond)
		close(imp.c)
	}()

	for isrc := range imp.c {
		i := isrc[0]
		src := isrc[1]

		if strings.HasPrefix(i, imp.Name) {
			if n, ok := imp.set[i]; ok {
				imp.set[i] = append(n, src)
				continue
			}
			p := strings.TrimPrefix(i, imp.Name)
			p = strings.Replace(p, "/", "", 1)
			go imp.collect(p)
		}
		n := imp.set[i]
		imp.set[i] = append(n, src)
	}
	return
}

func (imp *imports) collect(path string) {
	imp.wg.Add(1)
	defer imp.wg.Done()
	_ = fs.WalkDir(os.DirFS(imp.Dir), path, func(path2 string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == path2 {
				return nil
			}
			return fs.SkipDir
		}
		if strings.HasSuffix(path2, "_test.go") {
			return nil
		}
		src := filepath.Join(imp.Dir, path2)
		imps, err := Imports(src)
		if err != nil {
			return nil
		}
		for _, i := range imps {
			imp.c <- []string{i, src}
		}
		return nil
	})
}
