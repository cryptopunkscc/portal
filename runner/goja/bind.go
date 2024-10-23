package goja

import (
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/dop251/goja"
)

func Bind(vm *goja.Runtime, astral bind.Runtime) (err error) {
	var a = adapter{runtime: astral, vm: vm, queue: make(chan func(), 1024)}

	if err = vm.Set(bind.Log, a.Log); err != nil {
		return
	}
	if err = vm.Set(bind.Sleep, a.Sleep); err != nil {
		return
	}
	if err = vm.Set(bind.ServiceRegister, a.ServiceRegister); err != nil {
		return
	}
	if err = vm.Set(bind.ServiceClose, a.ServiceClose); err != nil {
		return
	}
	if err = vm.Set(bind.ConnAccept, a.ConnAccept); err != nil {
		return
	}
	if err = vm.Set(bind.ConnClose, a.ConnClose); err != nil {
		return
	}
	if err = vm.Set(bind.ConnWrite, a.ConnWrite); err != nil {
		return
	}
	if err = vm.Set(bind.ConnRead, a.ConnRead); err != nil {
		return
	}
	if err = vm.Set(bind.Query, a.Query); err != nil {
		return
	}
	if err = vm.Set(bind.QueryName, a.QueryName); err != nil {
		return
	}
	if err = vm.Set(bind.GetNodeInfo, a.NodeInfo); err != nil {
		return
	}
	if err = vm.Set(bind.ResolveId, a.Resolve); err != nil {
		return
	}
	if err = vm.Set(bind.Interrupt, a.Interrupt); err != nil {
		return
	}
	go func() {
		for f := range a.queue { // FIXME close queue on interrupt
			f()
		}
	}()
	return
}

type adapter struct {
	runtime bind.Runtime
	vm      *goja.Runtime
	queue   chan func()
}

func (a *adapter) Log(arg string) {
	a.runtime.Log(arg)
}
func (a *adapter) Sleep(millis int64) *goja.Promise {
	return a.promise0(func() { a.runtime.Sleep(millis) })
}
func (a *adapter) ServiceRegister(port string) *goja.Promise {
	return a.promise1(func() error { return a.runtime.ServiceRegister(port) })
}
func (a *adapter) ServiceClose(port string) *goja.Promise {
	return a.promise1(func() error { return a.runtime.ServiceClose(port) })
}
func (a *adapter) ConnAccept(port string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.runtime.ConnAccept(port) })
}
func (a *adapter) ConnClose(id string) *goja.Promise {
	return a.promise1(func() error { return a.runtime.ConnClose(id) })
}
func (a *adapter) ConnWrite(id string, data string) *goja.Promise {
	return a.promise1(func() error { return a.runtime.ConnWrite(id, data) })
}
func (a *adapter) ConnRead(id string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.runtime.ConnRead(id) })
}
func (a *adapter) Query(identity string, query string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.runtime.Query(identity, query) })
}
func (a *adapter) QueryName(name string, query string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.runtime.QueryName(name, query) })
}
func (a *adapter) Resolve(name string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.runtime.Resolve(name) })
}
func (a *adapter) NodeInfo(identity string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.runtime.NodeInfo(identity) })
}
func (a *adapter) Interrupt() *goja.Promise {
	return a.promise0(a.runtime.Interrupt)
}

func (a *adapter) promise0(f func()) *goja.Promise {
	return a.promise1(func() (err error) { f(); return })
}
func (a *adapter) promise1(f func() error) *goja.Promise {
	return a.promise2(func() (any, error) { return goja.Undefined(), f() })
}
func (a *adapter) promise2(f func() (any, error)) *goja.Promise {
	promise, resolve, reject := a.vm.NewPromise()
	go func() {
		val, err := f()
		a.queue <- func() {
			if err != nil {
				reject(err)
			} else {
				resolve(val)
			}
		}
	}()
	return promise
}
