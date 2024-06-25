package project

import "github.com/cryptopunkscc/portal/target"

var _ target.ProjectHtml = (*html)(nil)

type html struct {
	target.ProjectNpm
	target.Html
}

func (m *html) DistHtml() (t target.DistHtml) {
	return Dist[target.DistHtml](m)
}
