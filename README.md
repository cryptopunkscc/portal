# Portal

Desktop runtime & development environment for JavaScript decentralized applications driven by [Astral](https://github.com/cryptopunkscc/astrald/blob/master/docs/quickstart.md) network.

## Prerequisites

Before you can install & use Portal u need manually prepare the following dependencies:

### Astral

[Astral](https://github.com/cryptopunkscc/astrald/blob/master/docs/quickstart.md) is a core & mandatory networking dependency for Portal. It provides a plenty of features like p2p, encryption, identity, storage, and unified API for services and apps.

### Wails

Portal GUI runner uses Wails source code as a base, also it requires same dependencies for production and development purpose. For Installing them follow official Wails [installation](https://wails.io/docs/gettingstarted/installation) guideline.

### Clir

Command line interface for Portal is created using [clir](https://clir.leaanthony.com/) library.   

## Install

Portal sources can produce production or development executable. Generally speaking, the development runtime is an extended version of the production runtime.

* Install production runtime if you want to:
  * Run frontend application.
  * Run backend application.
* Install developer runtime if you want to:
  * Create new application.
  * Run development server.
  * Build application.
  * Generate application bundle.

### Linux & MacOS (Windows?)

Install production runtime: 
```shell
go install -tags "desktop,wv2runtime.download,production" github.com/cryptopunkscc/go-astral-js/cmd/portal
```

Install development runtime:
```shell
go install -tags dev github.com/cryptopunkscc/go-astral-js/cmd/portal
```

## How to use

Print help.

```shell
portal -help
```

Run development server.

```shell
portal dev ./example/wails
```

Create base application project from template.

```shell
portal create -n my_react_app -t react
```

Generate application bundle

```shell
portal bundle my_react_app
```

Run application bundle

```shell
portal open my_react_app/build/my_react_app.zip
```

### Legacy runners

v8 backend

```shell
go run ./cmd/legacy/v8 ./example/hello.js 
```

goja backend

```shell
go run ./cmd/legacy/goja ./example/hello.js 
```

WebView frontend

```shell
go run ./cmd/legacy/webview ./example/hello.html 
```

## Compatibility

Supported platforms for specific implementation.

* Frontend
    * wails - ES6
        * Linux
        * MacOS
        * Windows ?
    * WebView - ES6
        * Linux
        * MacOS ❌
        * Windows ?
* Backend
    * goja - ES6 (partial?)
        * Linux
        * MacOS
        * Windows ?
        * Android ?
    * V8 - ES6
        * Linux
        * MacOS ?
        * ~~Windows~~