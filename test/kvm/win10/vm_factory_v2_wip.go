package win10

import (
	"os"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
)

type VmFactory struct {
	VmName string

	OsDisk string
	OsIso  string

	VirtioIso string

	AutounattendIso string
	AutounattendDir string

	CidataDir string
	CidataIso string
}

func NewVmFactory() *VmFactory {
	return &VmFactory{
		VmName: "portal-test-win10",

		OsDisk: "./win.qcow2",
		OsIso:  "./iso/static/Win10_22H2_English_x64v1.iso",

		VirtioIso: "./iso/static/virtio-win-0.1.271.iso",

		AutounattendDir: "./iso/config/autounattend/",
		CidataDir:       "./iso/config/cidata/",

		AutounattendIso: "./iso/gen/autounattend.iso",
		CidataIso:       "./iso/gen/cidata.iso",
	}
}

func (vf *VmFactory) Init() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		*vf = *NewVmFactory()
	})
}

// VirtInstall installs and configures windows virtual machine.
func (vf *VmFactory) VirtInstall() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		execCmdRun(t, "virt-install",
			"--name", vf.VmName,
			"--ram", "8192",
			"--vcpus", "4",
			"--os-variant", "win10",
			"--network", "network=default",
			"--hvm",                                     // Hardware Virtual Machine. Makes OS unaware of being hosted in VM.
			"--disk", "path="+vf.OsDisk+",format=qcow2", // C:
			"--cdrom", vf.OsIso, // D:
			"--disk", "path="+vf.AutounattendIso+",device=cdrom", // E:
			"--disk", "path="+vf.VirtioIso+",device=cdrom", // F:
			"--disk", "path="+vf.CidataIso+",device=cdrom", // G:

			//"--graphics", "none", // Prevent displaying GUI.
			//"--console", "pty,target_type=serial",
		)
	},
		vf.GenerateAutounattendIso(),
		vf.GenerateCiDataIso(),
		vf.QemuImgCreate(),
	)
}

// GenerateAutounattendIso generates autounattend.iso that automates installation.
func (vf *VmFactory) GenerateAutounattendIso() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if fileExists(vf.AutounattendIso) {
			//return
		}
		execCmdRun(t,
			"genisoimage",
			"-output", vf.AutounattendIso,
			"-volid", "AUTOUNATTEND",
			"-joliet", "-rock",
			vf.AutounattendDir)
	})
}

// GenerateCiDataIso generates cidata.iso that contains NoCloud data like SSH key.
func (vf *VmFactory) GenerateCiDataIso() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if fileExists(vf.CidataIso) {
			//return
		}
		execCmdRun(t,
			"genisoimage",
			"-output", vf.CidataIso,
			"-volid", "cidata",
			"-joliet", "-rock",
			vf.CidataDir)
	})
}

// QemuImgCreate creates virtual drive for OS.
func (vf *VmFactory) QemuImgCreate() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if fileExists(vf.OsDisk) {
			return
		}
		execCmdRun(t, "qemu-img", "create", "-f", "qcow2", vf.OsDisk, "50G")
	},
		vf.QemuImgRemove(),
	)
}

func (vf *VmFactory) QemuImgRemove() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if !fileExists(vf.OsDisk) {
			return
		}
		err := os.Remove(vf.OsDisk)
		assert.NoError(t, err)
	},
		vf.VirtRemove(),
	)
}

func (vf *VmFactory) VirtRemove() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		if execCmd("virsh", "desc", vmName).Run() != nil {
			return // return if vm not exists
		}
		execCmdRun(t, "virsh", "destroy", vmName)
		execCmdRun(t, "virsh", "undefine", vmName)
	})
}
