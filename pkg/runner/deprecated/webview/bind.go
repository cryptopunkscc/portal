//go:build legacy

package webview

import (
	bind2 "github.com/cryptopunkscc/portal/pkg/bind"
	"github.com/webview/webview"
)

type WebView struct {
	webview.WebView
}

func New(debug bool) *WebView {
	return &WebView{WebView: webview.New(debug)}
}

func (view *WebView) BindApphost(core bind2.Core) {
	if err := view.Bind(bind2.Log, core.Log); err != nil {
		return
	}
	if err := view.Bind(bind2.ServiceRegister, core.ServiceRegister); err != nil {
		return
	}
	if err := view.Bind(bind2.ServiceClose, core.ServiceClose); err != nil {
		return
	}
	if err := view.Bind(bind2.ConnAccept, core.ConnAccept); err != nil {
		return
	}
	if err := view.Bind(bind2.ConnClose, core.ConnClose); err != nil {
		return
	}
	if err := view.Bind(bind2.ConnWrite, core.ConnWriteLn); err != nil {
		return
	}
	if err := view.Bind(bind2.ConnRead, core.ConnReadLn); err != nil {
		return
	}
	if err := view.Bind(bind2.Query, core.Query); err != nil {
		return
	}
	if err := view.Bind(bind2.GetNodeInfo, core.NodeInfo); err != nil {
		return
	}
	if err := view.Bind(bind2.ResolveId, core.Resolve); err != nil {
		return
	}
}
