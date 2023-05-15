package goja

import (
	astral_js "astral-js"
	"github.com/dop251/goja"
)

func Bind(vm *goja.Runtime, astral *astral_js.AppHostFlatAdapter) (err error) {
	var a = adapter{astral: astral, vm: vm}
	if err = vm.Set("log", astral.Log); err != nil {
		return
	}
	if err = vm.Set("sleep", a.Sleep); err != nil {
		return
	}
	if err = vm.Set("astral_port_listen", a.PortListen); err != nil {
		return
	}
	if err = vm.Set("astral_port_close", a.PortClose); err != nil {
		return
	}
	if err = vm.Set("astral_conn_accept", a.ConnAccept); err != nil {
		return
	}
	if err = vm.Set("astral_conn_close", a.ConnClose); err != nil {
		return
	}
	if err = vm.Set("astral_conn_write", a.ConnWrite); err != nil {
		return
	}
	if err = vm.Set("astral_conn_read", a.ConnRead); err != nil {
		return
	}
	if err = vm.Set("astral_dial", a.Dial); err != nil {
		return
	}
	if err = vm.Set("astral_dial_name", a.DialName); err != nil {
		return
	}
	if err = vm.Set("astral_node_info", a.NodeInfo); err != nil {
		return
	}
	if err = vm.Set("astral_resolve", a.Resolve); err != nil {
		return
	}
	return
}

type adapter struct {
	astral *astral_js.AppHostFlatAdapter
	vm     *goja.Runtime
}

func (a *adapter) Log(arg ...any) {
	a.astral.Log(arg...)
}

func (a *adapter) Sleep(millis int64) *goja.Promise {
	promise, resolve, _ := a.vm.NewPromise()
	go func() {
		a.astral.Sleep(millis)
		resolve(goja.Undefined())
	}()
	return promise
}

func (a *adapter) PortListen(port string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		err := a.astral.PortListen(port)
		if err != nil {
			reject(err)
		} else {
			resolve(goja.Undefined())
		}
	}()
	return promise
}

func (a *adapter) PortClose(port string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		err := a.astral.PortClose(port)
		if err != nil {
			reject(err)
		} else {
			resolve(goja.Undefined())
		}
	}()
	return promise
}

func (a *adapter) ConnAccept(port string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		conn, err := a.astral.ConnAccept(port)
		if err != nil {
			reject(err)
		} else {
			resolve(conn)
		}
	}()
	return promise
}

func (a *adapter) ConnClose(id string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		err := a.astral.ConnClose(id)
		if err != nil {
			reject(err)
		} else {
			resolve(goja.Undefined())
		}
	}()
	return promise
}

func (a *adapter) ConnWrite(id string, data string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		err := a.astral.ConnWrite(id, data)
		if err != nil {
			reject(err)
		} else {
			resolve(goja.Undefined())
		}
	}()
	return promise
}

func (a *adapter) ConnRead(id string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		val, err := a.astral.ConnRead(id)
		if err != nil {
			reject(err)
		} else {
			resolve(val)
		}
	}()
	return promise
}

func (a *adapter) Dial(identity string, query string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		val, err := a.astral.Dial(identity, query)
		if err != nil {
			reject(err)
		} else {
			resolve(val)
		}
	}()
	return promise
}

func (a *adapter) DialName(name string, query string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		val, err := a.astral.Dial(name, query)
		if err != nil {
			reject(err)
		} else {
			resolve(val)
		}
	}()
	return promise
}

func (a *adapter) Resolve(name string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		val, err := a.astral.Resolve(name)
		if err != nil {
			reject(err)
		} else {
			resolve(val)
		}
	}()
	return promise
}

func (a *adapter) NodeInfo(identity string) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		val, err := a.astral.NodeInfo(identity)
		if err != nil {
			reject(err)
		} else {
			resolve(val)
		}
	}()
	return promise
}
