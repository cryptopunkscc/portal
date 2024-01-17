# Astral JS - Golang

Astral JavaScript runtime environment written in golang for desktop.

## Platforms

Supported platforms for specific implementation:

* Frontend
  * WebView - ES6
    * Linux
    * MacOS ?
    * Windows ?
  * wails - ES6
    * Linux
    * MacOS ?
    * Windows ?
* Backend
  * V8 - ES6
    * Linux
    * MacOS ?
    * ~~Windows~~
  * goja - ES6 (partial?)
    * Linux
    * MacOS ?
    * Windows ?
    * Android ?

## Prerequisites

Make sure all required dependencies are installed.

* [WebView](https://github.com/webview/webview#prerequisites)
* [wails](https://wails.io/docs/gettingstarted/installation)

## Install

### Linux

Update binaries:

```shell
go build -o "$HOME/.local/bin/astral-runtime-webview" ./cmd/webview &&
go build -o "$HOME/.local/bin/astral-runtime-v8" ./cmd/v8 &&
go build -o "$HOME/.local/bin/astral-runtime-goja" ./cmd/goja
```

```shell
./cmd/wails/build.sh
```

Update anc

```shell
go build -o "$HOME/.local/bin/anc" github.com/cryptopunkscc/astrald/cmd/anc
```

Update config:

```shell
cp ./mod_apphost.yaml "$HOME/.config/astrald/config/"
```

## How to run

### AppHost

```shell
anc query localnode admin
```
start js app in admin console:
```
apphost run goja path_to_script.js
```

### Legacy

* v8 backend

```shell
go run ./cmd/v8 ./example/hello.js 
```

* goja backend

```shell
go run ./cmd/goja ./example/hello.js 
```

* WebView frontend

```shell
go run ./cmd/webview ./example/hello.html 
```

* wails frontend

```shell
./cmd/wails/build/bin/wails ./example/hello.html 
```

On MacOS:

```shell
./cmd/wails/build/bin/wails.app/Contents/MacOS/wails ./example/hello.html 
```
