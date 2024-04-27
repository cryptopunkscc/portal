package webview

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/webview/webview"
)

type WebView struct {
	webview.WebView
}

func New(debug bool) *WebView {
	return &WebView{WebView: webview.New(debug)}
}

func (view *WebView) BindApphost(astral apphost.Flat) {
	if err := view.Bind(apphost.Log, astral.Log); err != nil {
		return
	}
	if err := view.Bind(apphost.ServiceRegister, astral.ServiceRegister); err != nil {
		return
	}
	if err := view.Bind(apphost.ServiceClose, astral.ServiceClose); err != nil {
		return
	}
	if err := view.Bind(apphost.ConnAccept, astral.ConnAccept); err != nil {
		return
	}
	if err := view.Bind(apphost.ConnClose, astral.ConnClose); err != nil {
		return
	}
	if err := view.Bind(apphost.ConnWrite, astral.ConnWrite); err != nil {
		return
	}
	if err := view.Bind(apphost.ConnRead, astral.ConnRead); err != nil {
		return
	}
	if err := view.Bind(apphost.Query, astral.Query); err != nil {
		return
	}
	if err := view.Bind(apphost.QueryName, astral.QueryName); err != nil {
		return
	}
	if err := view.Bind(apphost.GetNodeInfo, astral.NodeInfo); err != nil {
		return
	}
	if err := view.Bind(apphost.Resolve, astral.Resolve); err != nil {
		return
	}
}
