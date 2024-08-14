# Developer Guide

This section explains how to start developing. Read if you want to create your own decentralized application(s) for
Astral network.

## Creating Project

Depending on your preferences and requirements you may need to develop one or many different applications that
can communicate through the network.

Portal provides predefined templates for generating empty projects. For listing available templates, execute:

```shell
portal-new -l
```

To create new project from template you can use the example command:

```shell
portal-new "html:gui js:service" ./example/echo
```

This command will generate 2 empty projects in the [`./example/echo`](./example/echo).
The fist is a html frontend named gui.
The second is a js backend named service.
For more advanced projects that require bundling, you can select different templates like `svelte` or `js-rollup`.

### Dependencies

For specific types of project you may need to install the following dependencies:

* [npm](https://nodejs.org/en/download/package-manager) - For npm based templates like `svelte` or `js-rollup`.
* [go](https://go.dev/doc/install) - For `go` project template.

## Application stages

During the development process portal application can exist in 3 different stages:

### Bundle

[Distributable](#Distributable) application source, packed as a single archive file contain manifest, ready for installation and run.

### Distributable

Distributable application directory, contains executable files, manifest, and assets ready to [bundle](#bundle).
Most often this is an intermediate stage between the project and the bundle, but for simple apps based on HTML or JS, it can be developed directly.

### Project

Contains not compiled application sources, assets, and manifest with optional build configuration.

## Developing Applications
 
Portal allows for developing and live testing HTML/JS applications using a hot-reloading server.

To run previously generated [echo](#creating-project) project, execute `portal-dev` command with a path to containing directory:

```shell
portal-dev ./example/echo
```

This command recursively searches a given path looking for [projects](#project), [distributables](#distributable) or [bundles](#bundle),
and executes each in its hot-reloading runner.
After running, besides the flood of logs in the console output, you should notice an opened empty application window.

Each time you save modified source code, the runner should reload the apps, which should be noticeable in logs and updated application UI.

## Communicating Applications

**Portal** provides embedded RPC library build upon the **Astral Apphost** interface.

To create simple echo server add the following snippet to [main.js](./example/echo/service/main.js).

```js
portal.rpc.serve({handlers: {echo: msg => msg + " " + msg}}).catch(portal.log)
```

This will register echo handler on dedicated **Astral Port** with name prefixed by application package.

Call echo handler by including the following script in [index.html](./example/echo/gui/index.html).

```html

<script>
    portal.rpc.call("new.portal.service.echo").request("hello astral").then(alert).catch(console.log)
</script>
```

which is equivalent of:

```html

<script>
    const service = portal.rpc.bind({"new.portal.service": ["echo"]})
    service.echo("hello astral").then(alert).catch(console.log)
</script>
```

After Saving modification:

1. The application window should reload automatically.
2. Frontend should send message to the service.
3. Service should respond with duplicated message.
4. Frontend should display response in popup.

NOTE: By default, headless apps close automatically with a 5-second timeout when there are no ongoing connections.
**Portal** can start closed applications on demand when another app tries to communicate.

If you don't want the app to close automatically, include a following code
in [portal.json](./example/echo/service/portal.json)

```json
{
  "env": {
    "timeout": -1
  }
}
```

To change the `timeout` length, specify the required amount in milliseconds.

If you need more examples of how to communicate **Astral Apps**, review projects in [example](./example) directory.

## Generating Application Bundle

To build application bundle ready to install on the user environment, use `portal-build <path>`, for example:

```shell
portal-build ./example/echo
```

This command will search the given directory for projects generating bundles into [build](./example/echo/build)
directory.
