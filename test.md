# Test commands

## Golang

Targets

```shell
go test ./runner/... ./resolve/...
```

RPC

```shell
go test ./pkg/rpc
```

## Cross integration

RPC

```shell
./make 8 && portal-dev -type 0 ./test/rpc
```

JS runtime lib

```shell
./make 8 && portal-dev -type 0 ./target/js/test
```
