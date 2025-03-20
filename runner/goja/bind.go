package goja

import (
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/dop251/goja"
)

func Bind(vm *goja.Runtime, core bind.Core) (err error) {
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	var a = adapter{core: core, vm: vm, queue: make(chan func(), 1024)}

	if err = vm.Set(bind.Exit, a.Exit); err != nil {
		return
	}
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
	if err = vm.Set(bind.ConnWriteLn, a.ConnWriteLn); err != nil {
		return
	}
	if err = vm.Set(bind.ConnReadLn, a.ConnReadLn); err != nil {
		return
	}
	if err = vm.Set(bind.Query, a.Query); err != nil {
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
	core  bind.Core
	vm    *goja.Runtime
	queue chan func()
}

func (a *adapter) Exit(code int) {
	a.core.Exit(code)
}
func (a *adapter) Log(arg string) {
	a.core.Log(arg)
}
func (a *adapter) Sleep(millis int64) *goja.Promise {
	return a.promise0(func() { a.core.Sleep(millis) })
}
func (a *adapter) ServiceRegister() *goja.Promise {
	return a.promise1(func() error { return a.core.ServiceRegister() })
}
func (a *adapter) ServiceClose() *goja.Promise {
	return a.promise1(func() error { return a.core.ServiceClose() })
}
func (a *adapter) ConnAccept() *goja.Promise {
	return a.promise2(func() (any, error) { return a.core.ConnAccept() })
}
func (a *adapter) ConnClose(id string) *goja.Promise {
	return a.promise1(func() error { return a.core.ConnClose(id) })
}
func (a *adapter) ConnWrite(id string, data []byte) *goja.Promise {
	return a.promise2(func() (any, error) { return a.core.ConnWrite(id, data) })
}
func (a *adapter) ConnRead(id string, n int) *goja.Promise {
	return a.promise2(func() (any, error) { return a.core.ConnRead(id, n) })
}
func (a *adapter) ConnWriteLn(id string, data string) *goja.Promise {
	return a.promise1(func() error { return a.core.ConnWriteLn(id, data) })
}
func (a *adapter) ConnReadLn(id string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.core.ConnReadLn(id) })
}
func (a *adapter) Query(target string, query string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.core.Query(target, query) })
}
func (a *adapter) Resolve(name string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.core.Resolve(name) })
}
func (a *adapter) NodeInfo(target string) *goja.Promise {
	return a.promise2(func() (any, error) { return a.core.NodeInfo(target) })
}
func (a *adapter) Interrupt() *goja.Promise {
	return a.promise0(a.core.Interrupt)
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
				reject(err.Error())
			} else {
				resolve(val)
			}
		}
	}()
	return promise
}
