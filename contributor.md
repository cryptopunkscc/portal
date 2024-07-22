# Contributor

This section explains how to build & install portal executables.
Read if you want to:

* generate `portal-installer`.
* test changes applied to portal source code.

## Required dependencies

Before you can start, make sure you have installed:

* go - https://go.dev/doc/install
* npm - https://nodejs.org/en/download/package-manager

## Building binaries

Portal repository includes an internal tool [`./cmd/make`](./cmd/make/make.go) for automating project build &
installation.

If you are starting from a fresh local copy of portal repository,
and want to verify the complete build & installation process, you can start just by executing:

```shell
./make
```

This shell script will trigger a building of required JS libs & embedded apps, install executables, and generate an
installer.

Optionally you can select specific steps by applying concatenated options.
For example, if you want to execute all steps except building this installer, execute:

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
<i|32> - Build portal installer into "./bin/portal-installer".
```

Is worth mentioning that the installation using:

* `<dp|24>` options depends on `go install` command so it outputs binaries into `$GOPATH/bin/`.
* `./bin/portal-installer` copies embedded binaries into `$HOME/local/bin/`.

