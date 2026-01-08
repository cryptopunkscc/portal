# portal-sdk

Development tools for Astral apps.

## Dependencies

Required for building portal-sdk from sources.

### Debian

```shell
sudo apt install gcc libgtk-3-dev libwebkit2gtk-4.1-dev
```

## Installation

```shell
go install -tags='dev,webkit2_41'
```

## Usage

```shell
$ portal-sdk help
```

### Creating a project

```shell
$ portal-sdk templates
```

This will list project templates.

```shell
$ portal-sdk new <template> <path>
```
This will generate a new application project from the specified template, using path base as the application name.

### Running project

```shell
$ portal-sdk run <path>
```

This will run the application project in the development runner with hot reloading.

### Building project

```shell
$ portal-sdk build [-pack] <path>
```

This will generate a distribution directory and, optionally, pack it into a bundle.
