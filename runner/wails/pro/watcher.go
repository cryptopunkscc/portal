package wails_pro

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/acarl005/stripansi"
	"golang.org/x/mod/semver"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// runViteWatcher will run the `frontend:dev:watcher` command if it was given, ex- `npm run dev`
func runViteWatcher(
	command string,
	frontendDir string,
	discoverViteURL bool,
) (close func(), viteUrl string, viteVersion string, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	scanner := newStdoutScanner()
	cmdSlice := strings.Split(command, " ")
	cmd := exec.CommandContext(ctx, cmdSlice[0], cmdSlice[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = scanner
	cmd.Dir = frontendDir
	setParentGID(cmd)

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, "", "", fmt.Errorf("unable to start frontend DevWatcher: %w", err)
	}

	if discoverViteURL {
		select {
		case serverURL := <-scanner.ViteServerURLChan:
			viteUrl = serverURL
		case <-time.After(time.Second * 10):
			cancel()
			return nil, "", "", errors.New("failed to find Vite server URL")
		}
	}

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

	close = func() {
		if atomic.CompareAndSwapInt32(&state, stateRunning, stateCanceling) {
			killProc(cmd, command)
		}
		cancel()
		wg.Wait()
	}
	return
}

// newStdoutScanner creates a new stdoutScanner
func newStdoutScanner() *stdoutScanner {
	return &stdoutScanner{
		ViteServerURLChan:  make(chan string, 2),
		ViteServerVersionC: make(chan string, 2),
	}
}

// stdoutScanner acts as a stdout target that will scan the incoming
// data to find out the vite server url
type stdoutScanner struct {
	ViteServerURLChan  chan string
	ViteServerVersionC chan string
	versionDetected    bool
}

// Write bytes to the scanner. Will copy the bytes to stdout
func (s *stdoutScanner) Write(data []byte) (n int, err error) {
	input := stripansi.Strip(string(data))
	if !s.versionDetected {
		v, err := detectViteVersion(input)
		if v != "" || err != nil {
			if err != nil {
				//logutils.LogRed("ViteStdoutScanner: %s", err)
				v = "v0.0.0"
			}
			s.ViteServerVersionC <- v
			s.versionDetected = true
		}
	}

	match := strings.Index(input, "Local:")
	if match != -1 {
		sc := bufio.NewScanner(strings.NewReader(input))
		for sc.Scan() {
			line := sc.Text()
			index := strings.Index(line, "Local:")
			if index == -1 || len(line) < 7 {
				continue
			}
			viteServerURL := strings.TrimSpace(line[index+6:])
			//logutils.LogGreen("Vite Server URL: %s", viteServerURL)
			_, err := url.Parse(viteServerURL)
			if err != nil {
				//logutils.LogRed(err.Error())
			} else {
				s.ViteServerURLChan <- viteServerURL
			}
		}
	}
	return os.Stdout.Write(data)
}

func detectViteVersion(line string) (string, error) {
	s := strings.Split(strings.TrimSpace(line), " ")
	if strings.ToLower(s[0]) != "vite" {
		return "", nil
	}

	if len(line) < 2 {
		return "", fmt.Errorf("unable to parse vite version")
	}

	v := s[1]
	if !semver.IsValid(v) {
		return "", fmt.Errorf("%s is not a valid vite version string", v)
	}

	return v, nil
}
