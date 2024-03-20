package publish

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	client "github.com/cryptopunkscc/astrald/lib/storage"
	"github.com/cryptopunkscc/astrald/mod/storage"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"io"
	"log"
	"os"
)

func Run(dir string) (err error) {
	r, err := runner.New(dir, runner.BundleTargets)
	if err != nil {
		return
	}

	c := client.NewClient(id.Anyone)
	targets := append(r.Backends, r.Frontends...)

	for _, t := range targets {
		log.Printf("publish %v", t.Path)
		if err = publish(c, t); err != nil {
			log.Printf("cannot publish %v: %v", t.Path, err)
		}
	}
	return
}

func publish(client *client.Client, target runner.Target) (err error) {
	dst, err := client.Create(&storage.CreateOpts{})
	if err != nil {
		return
	}
	src, err := os.Open(target.Path)
	if err != nil {
		return err
	}
	defer src.Close()
	l, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	log.Println("Commit", l, target.Path)
	dataID, err := dst.Commit()
	if err != nil {
		return
	}
	log.Printf("%v <- %v", dataID, target.Path)
	return
}
