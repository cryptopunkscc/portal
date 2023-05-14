package goja

import (
	astral_js "astral-js"
	"github.com/dop251/goja"
)

func Bind(vm *goja.Runtime, astral *astral_js.AppHostFlatAdapter) (err error) {
	if err = vm.Set("log", astral.Log); err != nil {
		return
	}
	if err = vm.Set("astral_port_listen", astral.PortListen); err != nil {
		return
	}
	if err = vm.Set("astral_port_close", astral.PortClose); err != nil {
		return
	}
	if err = vm.Set("astral_conn_accept", astral.ConnAccept); err != nil {
		return
	}
	if err = vm.Set("astral_conn_close", astral.ConnClose); err != nil {
		return
	}
	if err = vm.Set("astral_conn_write", astral.ConnWrite); err != nil {
		return
	}
	if err = vm.Set("astral_conn_read", astral.ConnRead); err != nil {
		return
	}
	if err = vm.Set("astral_dial", astral.Dial); err != nil {
		return
	}
	if err = vm.Set("astral_dial_name", astral.DialName); err != nil {
		return
	}
	if err = vm.Set("astral_node_info", astral.NodeInfo); err != nil {
		return
	}
	if err = vm.Set("astral_resolve", astral.Resolve); err != nil {
		return
	}
	return
}
