# Test commands

## Golang

Targets

```shell
go test ./resolve/... $(go list ./runner/... | grep -v /runner/webview)
```

RPC

```shell
go test ./runtime/rpc2
```

## Cross integration

RPC

```shell
./make 8 && portal-dev ./test/rpc
```

JS runtime lib

```shell
./make 8 && portal-dev ./runtime/js/test/common
```

```shell
./make 8 && portal-dev ./runtime/js/test/wails
```
