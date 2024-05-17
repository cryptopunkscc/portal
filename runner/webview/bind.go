package webview

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/webview/webview"
)

type WebView struct {
	webview.WebView
}

func New(debug bool) *WebView {
	return &WebView{WebView: webview.New(debug)}
}

func (view *WebView) BindApphost(astral target.Apphost) {
	if err := view.Bind(target.Log, astral.Log); err != nil {
		return
	}
	if err := view.Bind(target.ServiceRegister, astral.ServiceRegister); err != nil {
		return
	}
	if err := view.Bind(target.ServiceClose, astral.ServiceClose); err != nil {
		return
	}
	if err := view.Bind(target.ConnAccept, astral.ConnAccept); err != nil {
		return
	}
	if err := view.Bind(target.ConnClose, astral.ConnClose); err != nil {
		return
	}
	if err := view.Bind(target.ConnWrite, astral.ConnWrite); err != nil {
		return
	}
	if err := view.Bind(target.ConnRead, astral.ConnRead); err != nil {
		return
	}
	if err := view.Bind(target.Query, astral.Query); err != nil {
		return
	}
	if err := view.Bind(target.QueryName, astral.QueryName); err != nil {
		return
	}
	if err := view.Bind(target.GetNodeInfo, astral.NodeInfo); err != nil {
		return
	}
	if err := view.Bind(target.ResolveId, astral.Resolve); err != nil {
		return
	}
}
