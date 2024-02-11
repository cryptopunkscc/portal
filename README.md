# Portal

Astral apps JavaScript runtime environment for desktop.

## Prerequisites

To unlock connectivity features Portal requires [Astral](https://github.com/cryptopunkscc/astrald/blob/master/docs/quickstart.md) network installed, configured, and running.

## Install

Install dependencies from Wails [installation](https://wails.io/docs/gettingstarted/installation)
guide.

### Linux & MacOS

* Install production runtime if you want to:
    * Run frontend application.
    * Run backend application.

```shell
go install -tags "desktop,wv2runtime.download,production" github.com/cryptopunkscc/go-astral-js/cmd/portal
```

* Install developer runtime if you want to:
    * Create new application.
    * Run development server.
    * Build application.
    * Generate application bundle.

```shell
go install -tags dev github.com/cryptopunkscc/go-astral-js/cmd/portal
```

## How to use

* Print help.

```shell
portal -help
```

* Run development server.

```shell
portal dev ./example/wails
```

* Create base application project from template.

```shell
portal create -n my_react_app -t react
```

* Generate application bundle

```shell
portal bundle my_react_app
```

* Run application bundle

```shell
portal open my_react_app/build/my_react_app.zip
```

### Legacy commands

* v8 backend

```shell
go run ./cmd/legacy/v8 ./example/hello.js 
```

* goja backend

```shell
go run ./cmd/legacy/goja ./example/hello.js 
```

* WebView frontend

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