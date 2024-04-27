# Astral RPC

Basic extension to [lib astral](https://github.com/cryptopunkscc/astrald/tree/master/lib/astral). Provides a convenient wrapper for RPC protocol with customizable encoding.


## Usage

Implementing and running a service:

```go
package main

import rpc "github.com/cryptopunkscc/go-apphost-jrpc"

type service struct{}

func (s service) Sum(a int, b int) (c int, err error) {
	c = a + b
	return
}

func main() {
	err := rpc.NewApp("simple_calc").Run(service{})
	if err != nil {
		panic(err)   
	}
}
```

Calling method on the service:

```go
package main

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	rpc "github.com/cryptopunkscc/go-apphost-jrpc"
)

func main() {
	conn, _ := astral.Query(id.Identity{}, "simple_calc")
	r, _ := rpc.Query[int](conn, "sum", 2, 2)
	println(r)
}
```

See more comprehensive [example](./example).


## Protocol 

The general format of request is a command name followed byt arguments in a know format. 

example method:

```
<method name><arguments>
```

### Clir

The client can request data from Service by sending a method followed by commandline arguments.


```shell
methodName 1 true "string arg" -name "object arg"
```

### Json

The client can request data from Service by sending a method followed by positional arguments packed in array.

```json
methodName[1, true, "string arg", {"name": "object arg"}]
```

The service can respond by sending:
* `null` if there is nothing to send. 
* One error object.
* N amount of JSON objects.

example error:
```json
{"error": "some error message"}
```

The client can request a list of API methods provided by service by sending reserved method:

```json
["api"]
```

In response the service should send a JSON array containing method names:

```json
["api", "method1", "method2", "methodN"]
```
