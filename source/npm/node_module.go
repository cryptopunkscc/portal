package npm

import (
	"os"
	"os/exec"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
)

type NodeModule struct {
	source.Ref
	PackageJson PackageJson
}

func (p *NodeModule) ReadSrc(src source.Source) (err error) {
	return source.Readers{&p.Ref, &p.PackageJson}.ReadSrc(src)
}

func (p *NodeModule) WriteRef(ref source.Ref) (err error) {
	return source.Writers{&p.PackageJson, &p.Ref}.WriteRef(ref)
}

func (p *NodeModule) NpmInstall() (err error) {
	cmd := exec.Command("npm", "install")
	cmd.Dir = p.Path
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *NodeModule) Build() (err error) {
	if !p.PackageJson.CanBuild() {
		return plog.Errorf("missing scripts.build definition in %s/package.json", p.Path)
	}
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = p.Path
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
