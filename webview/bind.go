package webview

import (
	"astraljs"
	"github.com/webview/webview"
)

func Bind(view webview.WebView, astral *astraljs.AppHostFlatAdapter) {
	if err := view.Bind(astraljs.Log, astral.Log); err != nil {
		return
	}
	if err := view.Bind(astraljs.ServiceRegister, astral.ServiceRegister); err != nil {
		return
	}
	if err := view.Bind(astraljs.ServiceClose, astral.ServiceClose); err != nil {
		return
	}
	if err := view.Bind(astraljs.ConnAccept, astral.ConnAccept); err != nil {
		return
	}
	if err := view.Bind(astraljs.ConnClose, astral.ConnClose); err != nil {
		return
	}
	if err := view.Bind(astraljs.ConnWrite, astral.ConnWrite); err != nil {
		return
	}
	if err := view.Bind(astraljs.ConnRead, astral.ConnRead); err != nil {
		return
	}
	if err := view.Bind(astraljs.Query, astral.Query); err != nil {
		return
	}
	if err := view.Bind(astraljs.QueryName, astral.QueryName); err != nil {
		return
	}
	if err := view.Bind(astraljs.GetNodeInfo, astral.NodeInfo); err != nil {
		return
	}
	if err := view.Bind(astraljs.Resolve, astral.Resolve); err != nil {
		return
	}
}
