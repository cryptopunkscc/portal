package goja

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/dop251/goja"
)

func Bind(vm *goja.Runtime, astral target.Apphost) (err error) {
	var a = adapter{astral: astral, vm: vm, queue: make(chan func(), 1024)}

	if err = vm.Set(target.Log, a.Log); err != nil {
		return
	}
	if err = vm.Set(target.Sleep, a.Sleep); err != nil {
		return
	}
	if err = vm.Set(target.ServiceRegister, a.ServiceRegister); err != nil {
		return
	}
	if err = vm.Set(target.ServiceClose, a.ServiceClose); err != nil {
		return
	}
	if err = vm.Set(target.ConnAccept, a.ConnAccept); err != nil {
		return
	}
	if err = vm.Set(target.ConnClose, a.ConnClose); err != nil {
		return
	}
	if err = vm.Set(target.ConnWrite, a.ConnWrite); err != nil {
		return
	}
	if err = vm.Set(target.ConnRead, a.ConnRead); err != nil {
		return
	}
	if err = vm.Set(target.Query, a.Query); err != nil {
		return
	}
	if err = vm.Set(target.QueryName, a.QueryName); err != nil {
		return
	}
	if err = vm.Set(target.GetNodeInfo, a.NodeInfo); err != nil {
		return
	}
	if err = vm.Set(target.ResolveId, a.Resolve); err != nil {
		return
	}
	if err = vm.Set(target.Interrupt, a.Interrupt); err != nil {
		return
	}
	go func() {
		for f := range a.queue {
			f()
		}
	}()
	return
}

type adapter struct {
	astral target.Apphost
	vm     *goja.Runtime
	queue  chan func()
}

func (a *adapter) Log(arg ...any) {
	a.astral.LogArr(arg)
}
func (a *adapter) Sleep(millis int64) *goja.Promise {
	return a.promise0(func() { a.astral.Sleep(millis) })
}
func (a *adapter) ServiceRegister(port string) *goja.Promise {
	return a.promise1(func() error { return a.astral.ServiceRegister(port) })
}
func (a *adapter) ServiceClose(port string) *goja.Promise {
	return a.promise1(func() error { return a.astral.ServiceClose(port) })
}
func (a *adapter) ConnAccept(port string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.astral.ConnAccept(port) })
}
func (a *adapter) ConnClose(id string) *goja.Promise {
	return a.promise1(func() error { return a.astral.ConnClose(id) })
}
func (a *adapter) ConnWrite(id string, data string) *goja.Promise {
	return a.promise1(func() error { return a.astral.ConnWrite(id, data) })
}
func (a *adapter) ConnRead(id string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.astral.ConnRead(id) })
}
func (a *adapter) Query(identity string, query string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.astral.Query(identity, query) })
}
func (a *adapter) QueryName(name string, query string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.astral.QueryName(name, query) })
}
func (a *adapter) Resolve(name string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.astral.Resolve(name) })
}
func (a *adapter) NodeInfo(identity string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.astral.NodeInfo(identity) })
}
func (a *adapter) Interrupt() *goja.Promise {
	return a.promise0(a.astral.Interrupt)
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
