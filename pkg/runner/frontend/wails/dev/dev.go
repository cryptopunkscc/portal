package dev

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// runFrontendDevWatcherCommand will run the `frontend:dev:watcher` command if it was given, ex- `npm run dev`
func runFrontendDevWatcherCommand(frontendDirectory string, devCommand string, discoverViteServerURL bool) (func(), string, string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	scanner := NewStdoutScanner()
	cmdSlice := strings.Split(devCommand, " ")
	cmd := exec.CommandContext(ctx, cmdSlice[0], cmdSlice[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = scanner
	cmd.Dir = frontendDirectory
	setParentGID(cmd)

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, "", "", fmt.Errorf("unable to start frontend DevWatcher: %w", err)
	}

	var viteServerURL string
	if discoverViteServerURL {
		select {
		case serverURL := <-scanner.ViteServerURLChan:
			viteServerURL = serverURL
		case <-time.After(time.Second * 10):
			cancel()
			return nil, "", "", errors.New("failed to find Vite server URL")
		}
	}

	viteVersion := ""
	select {
	case version := <-scanner.ViteServerVersionC:
		viteVersion = version

	case <-time.After(time.Second * 5):
		// That's fine, then most probably it was not vite that was running
	}

	//logutils.LogGreen("Running frontend DevWatcher command: '%s'", devCommand)
	var wg sync.WaitGroup
	wg.Add(1)

	const (
		stateRunning   int32 = 0
		stateCanceling int32 = 1
		stateStopped   int32 = 2
	)
	state := stateRunning
	go func() {
		if err := cmd.Wait(); err != nil {
			wasRunning := atomic.CompareAndSwapInt32(&state, stateRunning, stateStopped)
			if err.Error() != "exit status 1" && wasRunning {
				//logutils.LogRed("Error from DevWatcher '%s': %s", devCommand, err.Error())
			}
		}
		atomic.StoreInt32(&state, stateStopped)
		wg.Done()
	}()

	return func() {
		if atomic.CompareAndSwapInt32(&state, stateRunning, stateCanceling) {
			killProc(cmd, devCommand)
		}
		cancel()
		wg.Wait()
	}, viteServerURL, viteVersion, nil
}
