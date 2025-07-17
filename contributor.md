# Contributor

This section explains how to build & install portal executables. Read if you want to:

* generate `install-portal-to-astral`.
* test changes applied to portal source code.

## Required dependencies

Before you can start, make sure you have installed:

* go - https://go.dev/doc/install
* npm - https://nodejs.org/en/download/package-manager

## Building binaries

Portal repository includes an internal tool [`./cmd/make`](./cmd/make/make.go) for automating project build & installation.

If you are starting from a fresh local copy of portal repository, 
and want to verify the complete build & installation process, execute from the repository root:

```shell
./make
```

This shell script is a shortcut to `./cmd/make`.
When called without arguments it will run all building steps.

Optionally you can select specific steps by applying concatenated options.
For example, if you want to execute all steps except building `install-portal-to-astral`, execute:

```shell
./make ladp
```

which is the same as:

```shell
./make 30
```

Here is a complete list of modifiers for make tool:

```
<empty> - Select all options.
<l|2> - Build JS libraries. (required by embed JS apps)
<a|4> - Build embed JS applications. (required by portal binaries)
<d|8> - Install portal for developer binary into "$HOME/go/bin/portal-dev".
<p|16> - Install portal for user into "$HOME/go/bin/portal". 
<i|32> - Build portal installer into "./bin/install-portal-to-astral".
```

Is worth mentioning that the installation run using:

* `./make <dp|24>` depends on `go install` command so it outputs binaries into `$GOPATH/bin/`.
* `./bin/install-portal-to-astral` copies embedded binaries into platform specific directory.

## GOOS

To build `install-portal-to-astral` for specific platforms provide them as additional arguments:

```shell
./make i linux windows
```