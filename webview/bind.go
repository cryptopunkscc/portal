package webview

import (
	"astral-js"
	"github.com/webview/webview"
)

func Bind(view webview.WebView, astral *astral_js.AppHostFlatAdapter) {
	if err := view.Bind("log", astral.Log); err != nil {
		return
	}
	if err := view.Bind("astral_port_listen", astral.PortListen); err != nil {
		return
	}
	if err := view.Bind("astral_port_close", astral.PortClose); err != nil {
		return
	}
	if err := view.Bind("astral_conn_accept", astral.ConnAccept); err != nil {
		return
	}
	if err := view.Bind("astral_conn_close", astral.ConnClose); err != nil {
		return
	}
	if err := view.Bind("astral_conn_write", astral.ConnWrite); err != nil {
		return
	}
	if err := view.Bind("astral_conn_read", astral.ConnRead); err != nil {
		return
	}
	if err := view.Bind("astral_dial", astral.Dial); err != nil {
		return
	}
	if err := view.Bind("astral_dial_name", astral.DialName); err != nil {
		return
	}
	if err := view.Bind("astral_node_info", astral.NodeInfo); err != nil {
		return
	}
	if err := view.Bind("astral_resolve", astral.Resolve); err != nil {
		return
	}
}
