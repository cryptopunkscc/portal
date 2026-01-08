# portal-wails

Runs Astral HTML apps using [wails](https://wails.io/)-based [runner](../../runner/v2/wails).

## Dependencies

Required to build from sources:

### Debian

```shell
sudo apt install gcc libgtk-3-dev libwebkit2gtk-4.1-dev
```

## Installation

```shell
go install -tags='desktop,wv2runtime.download,production,webkit2_41' -ldflags='-w -s'
```

Or

```shell
go install -tags='dev,webkit2_41'
```

To enable wails development features like: webview inspect mode, dev server, etc...

## Usage

```shell
$ portal-wails <argument> 
```

where `<argument>` can be:

* **name** declared in the application's portal.json (only if app was published)
* **package** declared in the application's portal.json (only if app was published)
* local **path** to the application distribution directory
* local path to the application **bundle**
