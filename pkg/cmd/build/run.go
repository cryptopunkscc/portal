package build

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"log"
	"os"
	"os/exec"
	"sync"
)

func Run(dir string) (err error) {
	r, err := runner.New(dir, runner.DevTargets)
	if err != nil {
		return
	}
	return RunRunner(r)
}

func RunRunner(r *runner.Runner) (err error) {
	targets := append(r.Backends, r.Frontends...)

	wait := sync.WaitGroup{}
	wait.Add(len(targets))
	for _, target := range targets {
		go func(src string) {
			defer wait.Done()
			if err := npmInstall(src); err != nil {
				log.Println(err)
			}
			if err := npmRunBuild(src); err != nil {
				log.Println(err)
			}
		}(target.Path)
	}
	wait.Wait()
	return
}

func npmInstall(dir string) error {
	return runCmd(dir, "npm", "install")
}

func npmRunBuild(dir string) error {
	return runCmd(dir, "npm", "run", "build")
}

func runCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
