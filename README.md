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
* goja - ES6 (partial?)
  * Linux 
  * MacOS ?
  * Windows ?
  * Android ?

## Prerequisites

Make sure all required dependencies are installed.

* [WebView](https://github.com/webview/webview#prerequisites)

## How to run

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
