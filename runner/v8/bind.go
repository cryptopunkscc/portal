package v8

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"log"
	v8 "rogchap.com/v8go"
)

func Bind(iso *v8.Isolate, astral target.Apphost) (template *v8.ObjectTemplate, err error) {
	template = v8.NewObjectTemplate(iso)
	a := adapter{astral}
	if err = template.Set(target.Log, v8.NewFunctionTemplate(iso, a.Log)); err != nil {
		return
	}
	if err = template.Set(target.Sleep, v8.NewFunctionTemplate(iso, a.Sleep)); err != nil {
		return
	}
	if err = template.Set(target.ServiceRegister, v8.NewFunctionTemplate(iso, a.ServiceRegister)); err != nil {
		return
	}
	if err = template.Set(target.ServiceClose, v8.NewFunctionTemplate(iso, a.ServiceClose)); err != nil {
		return
	}
	if err = template.Set(target.ConnAccept, v8.NewFunctionTemplate(iso, a.ConnAccept)); err != nil {
		return
	}
	if err = template.Set(target.ConnClose, v8.NewFunctionTemplate(iso, a.ConnClose)); err != nil {
		return
	}
	if err = template.Set(target.ConnWrite, v8.NewFunctionTemplate(iso, a.ConnWrite)); err != nil {
		return
	}
	if err = template.Set(target.ConnRead, v8.NewFunctionTemplate(iso, a.ConnRead)); err != nil {
		return
	}
	if err = template.Set(target.Query, v8.NewFunctionTemplate(iso, a.Query)); err != nil {
		return
	}
	if err = template.Set(target.QueryName, v8.NewFunctionTemplate(iso, a.QueryName)); err != nil {
		return
	}
	if err = template.Set(target.GetNodeInfo, v8.NewFunctionTemplate(iso, a.NodeInfo)); err != nil {
		return
	}
	if err = template.Set(target.ResolveId, v8.NewFunctionTemplate(iso, a.Resolve)); err != nil {
		return
	}
	return
}

type adapter struct {
	astral target.Apphost
}

func (a *adapter) Log(info *v8.FunctionCallbackInfo) *v8.Value {
	log.Println(info.Args())
	return nil
}

func (a *adapter) Sleep(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	t := info.Args()[0].Integer()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		a.astral.Sleep(t)
		resolver.Resolve(v8.Undefined(iso))
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) ServiceRegister(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	port := info.Args()[0].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		err := a.astral.ServiceRegister(port)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
		}
		resolver.Resolve(v8.Undefined(iso))
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) ServiceClose(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	port := info.Args()[0].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		err := a.astral.ServiceClose(port)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
		}
		resolver.Resolve(v8.Undefined(iso))
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) ConnAccept(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	port := info.Args()[0].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		id, err := a.astral.ConnAccept(port)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
			return
		}
		val, err := v8.NewValue(iso, id)
		if err != nil {
			log.Fatal(err)
		}
		resolver.Resolve(val)
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) ConnClose(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	id := info.Args()[0].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		err := a.astral.ConnClose(id)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
		}
		resolver.Resolve(v8.Undefined(iso))
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) ConnWrite(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	id := info.Args()[0].String()
	data := info.Args()[1].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		err := a.astral.ConnWrite(id, data)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
		}
		resolver.Resolve(v8.Undefined(iso))
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) ConnRead(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	id := info.Args()[0].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		data, err := a.astral.ConnRead(id)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
			return
		}
		val, err := v8.NewValue(iso, data)
		if err != nil {
			log.Fatal(err)
		}
		resolver.Resolve(val)
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) Query(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	id := info.Args()[0].String()
	query := info.Args()[1].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		connId, err := a.astral.Query(id, query)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
			return
		}
		val, err := v8.NewValue(iso, connId)
		if err != nil {
			log.Fatal(err)
		}
		resolver.Resolve(val)
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) QueryName(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	name := info.Args()[0].String()
	query := info.Args()[1].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		connId, err := a.astral.QueryName(name, query)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
			return
		}
		val, err := v8.NewValue(iso, connId)
		if err != nil {
			log.Fatal(err)
		}
		resolver.Resolve(val)
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) NodeInfo(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	id := info.Args()[0].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		nodeInfo, err := a.astral.NodeInfo(id)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
			return
		}
		val, err := v8.NewValue(iso, nodeInfo)
		if err != nil {
			log.Fatal(err)
		}
		resolver.Resolve(val)
	}()
	return resolver.GetPromise().Value
}

func (a *adapter) Resolve(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	name := info.Args()[0].String()
	resolver, _ := v8.NewPromiseResolver(info.Context())
	go func() {
		nodeInfo, err := a.astral.Resolve(name)
		if err != nil {
			val, err := v8.NewValue(iso, err.Error())
			if err != nil {
				log.Fatal(err)
			}
			resolver.Reject(val)
			return
		}
		val, err := v8.NewValue(iso, nodeInfo)
		if err != nil {
			log.Fatal(err)
		}
		resolver.Resolve(val)
	}()
	return resolver.GetPromise().Value
}
