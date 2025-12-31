package source

import (
	"github.com/spf13/afero"
)

type Provider interface {
	GetSource(source string) Source
}

type Providers []Provider

func (p Providers) GetSource(source string) (out Source) {
	for _, provider := range p {
		if out = provider.GetSource(source); out != nil {
			return
		}
	}
	return
}

type PathProvider struct{ Ref }

func (r PathProvider) GetSource(path string) (out Source) {
	if err := r.Ref.Checkout(path); err == nil {
		out = r.New()
	}
	return
}

var OsFs = PathProvider{Ref{Fs: afero.NewOsFs()}}
