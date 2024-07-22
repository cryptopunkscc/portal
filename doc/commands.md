### Example commands

Print help.

```shell
portal --help
```

Run applications from directory.

```shell
portal ../example/rpc
```

Run applications in development server supporting hot reload and node module based projects.

```shell
portal-dev ../example/project
```

Build & generate application bundles.

```shell
portal-dev b ../example/project
```

Install applications from generated bundles.

```shell
portal-app i ../example/project
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
portal-dev c "html:frontend js:backend" ../my_project
```

List available templates.

```shell
portal-dev c -l
```
