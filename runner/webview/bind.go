//go:build legacy

package webview

import (
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/webview/webview"
)

type WebView struct {
	webview.WebView
}

func New(debug bool) *WebView {
	return &WebView{WebView: webview.New(debug)}
}

func (view *WebView) BindApphost(core bind.Core) {
	if err := view.Bind(bind.Log, core.Log); err != nil {
		return
	}
	if err := view.Bind(bind.ServiceRegister, core.ServiceRegister); err != nil {
		return
	}
	if err := view.Bind(bind.ServiceClose, core.ServiceClose); err != nil {
		return
	}
	if err := view.Bind(bind.ConnAccept, core.ConnAccept); err != nil {
		return
	}
	if err := view.Bind(bind.ConnClose, core.ConnClose); err != nil {
		return
	}
	if err := view.Bind(bind.ConnWrite, core.ConnWriteLn); err != nil {
		return
	}
	if err := view.Bind(bind.ConnRead, core.ConnReadLn); err != nil {
		return
	}
	if err := view.Bind(bind.Query, core.Query); err != nil {
		return
	}
	if err := view.Bind(bind.GetNodeInfo, core.NodeInfo); err != nil {
		return
	}
	if err := view.Bind(bind.ResolveId, core.Resolve); err != nil {
		return
	}
}
