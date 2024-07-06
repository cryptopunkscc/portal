# Portal

Runtime & development environment for decentralized applications driven by [Astral](https://github.com/cryptopunkscc/astrald/blob/master/docs/quickstart.md) network.

## Installation

To build portal executables run:

```shell
./make <number?>
```

Where optional <number?> can be concatenation of the following options :

```
<empty> - Select all options.
<1> - Install dependencies.
<2> - Build JS libraries.
<4> - Build embed JS applications.
<8> - Install portal for developer into "$HOME/go/bin/portal-dev".
<16> - Install portal for user into "$HOME/go/bin/portal". 
<32> - Build portal installer.
```

See examples: [commands.md](./test.md)

### Example commands

Print help.

```shell
portal -help
```

Run applications from directory.

```shell
portal ./example/rpc
```

Run applications in development server supporting hot reload and node module based projects.

```shell
portal-dev ./example/project
```

Build & generate application bundles.

```shell
portal-dev b ./example/project
```

Install applications from generated bundles.

```shell
portal i ./example/project
```

Run application by name.

```shell
portal launcher
```

Run application by package name.

```shell
portal example.project.svelte
```

Create new project from template.

```shell
portal-dev c "html:frontend js:backend" ./my_project
```

List available templates.

```shell
portal-dev c -l
```
