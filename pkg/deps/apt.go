package deps

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
)

func AptInstallMissing(deps []string) (err error) {
	var missing []string
	for _, d := range deps {

		if err := exec.Call(".", "dpkg-query", "-l", d); err != nil {
			log.Printf("missing dep: %s, %v", d, err)
			missing = append(missing, d)
		}
	}
	if len(missing) > 0 {
		cmd := append([]string{"sudo", "apt", "install"}, missing...)
		if err = exec.Run(".", cmd...); err != nil {
			panic(err)
			return err
		}
	}
	return
}

func Check(cmd ...string) (err error) {
	if err = exec.Call(".", cmd...); err != nil {
		err = fmt.Errorf("%s is required but not installed", cmd[0])
	}
	return
}
