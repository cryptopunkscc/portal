package astrald

func (i *Initializer) initNodeResources() (err error) {
	if i.resources.Path == "" {
		i.resources.Path = i.NodeRoot
		err = i.resources.Init()
	}
	return
}
