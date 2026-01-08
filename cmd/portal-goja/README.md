# portal-goja

Runs Astral JS services using [goja](https://github.com/dop251/goja)-based [runner](../../runner/v2/goja).

## Installation

```shell
go install
```

## Usage

```shell
$ portal-goja <argument> 
```
where `<argument>` can be:
* **name** declared in the application's portal.json (only if app was published)
* **package** declared in the application's portal.json (only if app was published)
* local **path** to the application distribution directory
* local path to the application **bundle**
