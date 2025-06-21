package test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestE2E(t *testing.T) {
	image := "e2e-test"
	container1 := "e2e-test-1"

	tests := []struct {
		Name, Dir string
		Cmd       []string
	}{
		{Name: "build portal installer", Dir: "../", Cmd: []string{
			"./mage", "build:out", "./test", "build:installer"}},
		{Name: "build_image", Cmd: []string{
			"docker", "build", "-t", image + ":latest", "."}},
		{Name: "run_container1", Cmd: []string{
			"docker", "run", "--name", container1, "-p", "8081:8080", image}},
		//"docker", "run", "-d", "--name", container1, "--network", network, "-p", "8081:8080", image}}, // FIXME custom network doesn't work for some reason
		{Name: "stop_container1", Cmd: []string{
			"docker", "stop", container1 /*, instance2*/}},
		{Name: "remove_container1", Cmd: []string{
			"docker", "rm", container1 /*, instance2*/}},
		{Name: "remove_image", Cmd: []string{
			"docker", "rmi", image}},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.NotEmpty(t, test.Cmd)
			c := exec.Command(test.Cmd[0], test.Cmd[1:]...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			c.Env = append(os.Environ())
			c.Dir = "."
			if test.Dir != "" {
				c.Dir = test.Dir
			}
			err := c.Run()
			assert.NoError(t, err)
		})
	}
}
