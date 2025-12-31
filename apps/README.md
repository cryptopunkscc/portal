# Apps (DEPRECATED)

The following directory contains default apps embedded in the portal environment.

## Running apps without portald

The portald daemon automates app management. However, it isn't mandatory to run astral apps.

### Building default apps

Build the default apps:

```shell
./mage build:apps
```

optionally, you can build default apps by running [build_test.go](./build_test.go)

### Running app binary

By default, the app binary is located at `{app}/dist/{os}/{arch}/main`.

You can run binary directly from the command-line, for example:

```shell
./apps/html-dev/dist/linux/amd64/main help
```

### Development

Create a new svelte app project:

```shell
./apps/html-dev/dist/linux/amd64/main new -t svelte ./my-svelte-app
```

Run the app directly from a project using the development runner to take advantage of hot-reloading:

```shell
ASTRALD_APPHOST_TOKEN=<token> ./apps/html-dev/dist/linux/amd64/main ./my-svelte-app
```

Create app bundle:

```shell
./apps/builder/dist/linux/amd64/main -pack ./my-svelte-app
```

Run the app bundle using the production runner:

```shell
ASTRALD_APPHOST_TOKEN=<token> ./apps/html/dist/linux/amd64/main ./my-svelte-app/build/my.app.my-svelte-app_0.0.0.portal
```
