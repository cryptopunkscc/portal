package win10

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
)

const (
	vmName  = "portal-test-win10"
	winDisk = "./win.qcow2"

	winIso    = "./iso/static/Win10_22H2_English_x64v1.iso"
	virtioIso = "./iso/static/virtio-win-0.1.271.iso"

	autounattendDir = "./iso/config/autounattend/"
	cidataDir       = "./iso/config/cidata/"
	sshKeyPub       = "./iso/config/cidata/sshkey.pub"

	autounattendIso = "./iso/gen/autounattend.iso"
	cidataIso       = "./iso/gen/cidata.iso"
)

var (
	winSshKey    = ".ssh/" + vmName
	winSshKeyPub = ".ssh/" + vmName + ".pub"
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	winSshKey = filepath.Join(home, winSshKey)
	winSshKeyPub = filepath.Join(home, winSshKeyPub)
}

func removeVmDisk() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if !fileExists(winDisk) {
			return
		}
		err := os.Remove(winDisk)
		assert.NoError(t, err)
	})
}

func virtUndefine() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if execCmd("virsh", "desc", vmName).Run() != nil {
			return // return if vm not exists
		}
		execCmdRun(t, "virsh", "destroy", vmName)
		execCmdRun(t, "virsh", "undefine", vmName)
	})
}

// virtInstall installs and configures windows virtual machine.
func virtInstall() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		execCmdRun(t, "virt-install",
			"--name", vmName,
			"--ram", "8192",
			"--vcpus", "4",
			"--os-variant", "win10",
			"--network", "network=default",
			//"--graphics", "none", // Prevent displaying GUI.
			"--hvm",                                   // Hardware Virtual Machine. Makes OS unaware of being hosted in VM.
			"--disk", "path="+winDisk+",format=qcow2", // C:
			"--cdrom", winIso, // D:
			"--disk", "path="+autounattendIso+",device=cdrom", // E:
			"--disk", "path="+virtioIso+",device=cdrom", // F:
			"--disk", "path="+cidataIso+",device=cdrom", // G:
			//"--disk", "path="+cloudbaseIso+",device=cdrom", // H:
			//"--console", "pty,target_type=serial",
		)
	})
}

// genCiDataIso generates cidata.iso that contains NoCloud data like SSH key.
func genCiDataIso() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if fileExists(cidataIso) {
			//return
		}
		execCmdRun(t,
			"genisoimage",
			"-output", cidataIso,
			"-volid", "cidata",
			"-joliet", "-rock",
			cidataDir)
	})
}

// genAutounattendIso generates autounattend.iso that automates installation.
func genAutounattendIso() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if fileExists(autounattendIso) {
			//return
		}
		execCmdRun(t,
			"genisoimage",
			"-output", autounattendIso,
			"-volid", "AUTOUNATTEND",
			"-joliet", "-rock",
			autounattendDir)
	})
}

// createQemuImage creates virtual drive for OS.
func createQemuImage() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if fileExists(winDisk) {
			return
		}
		execCmdRun(t, "qemu-img", "create", "-f", "qcow2", winDisk, "50G")
	})
}

func copyPubKey() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		execCmdRun(t, "cp", winSshKeyPub, sshKeyPub)
	})
}

// genSshKey generates SSH key that will be injected to VM.
func genSshKey() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if fileExists(winSshKey) {
			return
		}
		execCmdRun(t, "ssh-keygen", "-t", "rsa", "-b", "4096", "-f", winSshKey)
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
