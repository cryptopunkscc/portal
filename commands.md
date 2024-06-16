## Install

Install portal, portal-dev, and required dependencies.

```shell
go run ./cmd/install
```

Run portal 

## Tests

### Golang

Targets

```shell
go test ./test/target/
```

RPC

```shell
go test ./pkg/rpc
```

### Cross integration

RPC

```shell
go run ./cmd/install 8 && portal-dev -type 0 ./test/rpc
```

JS runtime lib

```shell
go run ./cmd/install 8 && portal-dev -type 0 ./target/js/test
```
