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

func (view *WebView) BindApphost(runtime bind.Runtime) {
	if err := view.Bind(bind.Log, runtime.Log); err != nil {
		return
	}
	if err := view.Bind(bind.ServiceRegister, runtime.ServiceRegister); err != nil {
		return
	}
	if err := view.Bind(bind.ServiceClose, runtime.ServiceClose); err != nil {
		return
	}
	if err := view.Bind(bind.ConnAccept, runtime.ConnAccept); err != nil {
		return
	}
	if err := view.Bind(bind.ConnClose, runtime.ConnClose); err != nil {
		return
	}
	if err := view.Bind(bind.ConnWrite, runtime.ConnWriteLn); err != nil {
		return
	}
	if err := view.Bind(bind.ConnRead, runtime.ConnReadLn); err != nil {
		return
	}
	if err := view.Bind(bind.Query, runtime.Query); err != nil {
		return
	}
	if err := view.Bind(bind.QueryName, runtime.QueryName); err != nil {
		return
	}
	if err := view.Bind(bind.GetNodeInfo, runtime.NodeInfo); err != nil {
		return
	}
	if err := view.Bind(bind.ResolveId, runtime.Resolve); err != nil {
		return
	}
}
