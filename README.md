# Astral JS - Golang

Astral JavaScript runtime environment written in golang for desktop.

## Platforms

Supported platforms for specific implementation:

* WebView - ES6
  * Linux
  * MacOS ?
  * Windows ?
* V8 - ES6
  * Linux
  * MacOS ?
  * ~~Windows~~
* goja - ES5
  * Linux 
  * MacOS ?
  * Windows ?
  * Android ?

## Prerequisites

Make sure all required dependencies are installed.

* [WebView](https://github.com/webview/webview#prerequisites)

## How to run

* v8 backend + es6
```shell
go run ./cmd/v8 ./example/hello.js 
```

* v8 backend (es5)
```shell
go run ./cmd/v8 ./example/hello.es5.js 
```

* goja + es5
```shell
go run ./cmd/goja ./example/hello.es5.js 
```

Start frontend:

```shell
go run ./cmd/webview ./example/hello.html 
```
