package setup

import apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"

func (r *Runner) initApphostConfig() {
	if r.ApphostConfig == nil {
		if err := r.readApphostConfig(); err != nil {
			r.ApphostConfig = &apphost.Config{}
			return
		}
		r.log.Println("loaded existing apphost config")
	}
}

func (r *Runner) readApphostConfig() (err error) {
	return r.resources.ReadYaml(apphostYaml, &r.ApphostConfig)
}
func (r *Runner) writeApphostConfig() (err error) {
	return r.resources.WriteYaml(apphostYaml, r.ApphostConfig)
}

const apphostYaml = "apphost.yaml"
