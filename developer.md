# Developer Guide

This section explains how to start developing. Read if you want to create your own decentralized application(s) for
Astral network.

## Creating Project

Depending on your preferences and requirements you may need to develop one or many different applications that
can communicate through the network.
You can choose from predefined templates for generating empty projects. For listing available
templates, execute:

```shell
portal-dev c -l
```

To create new project from template you can use the example command:

```shell
portal-dev c "html:gui js:service" ./example/echo-portal
```

This command will generate 2 empty projects in the [`./example/echo-portal`](./example/echo-portal).
The fist is a html frontend named gui.
The second is a js backend named service.
For more advanced projects that require bundling, you can select different templates like `sevlte` or `js-rollup`.

### Dependencies

For specific types of project you may need to install the following dependencies:

* [npm](https://nodejs.org/en/download/package-manager) - For npm based templates like `svelte` or `js-rollup`.
* [go](https://go.dev/doc/install) - For `go` project template.

## Developing Applications

In general, you can just write your code and generate a bundle,
but a more convenient way to develop JS applications is live testing the changes using hot-reloading server.

To run your projects execute `portal-dev` command with a path to containing directory as argument, like:

```shell
portal-dev ./example/echo-portal
```

This command recursively searches a given path looking for projects or apps,
then executes each in its hot-reloading runner.
After running, besides the flood of logs in the console output, you should notice an empty window opened.

Each time you save changes in the project source code, the runner should reload the apps. That should be noticeable in
logs and updated application UI.

## Communicating Applications

**Portal** provides embedded RPC library build upon the **Astral Apphost** interface. 

For example to create simple echo server add the following snippet
to [main.js](./example/echo-portal/service/main.js).

```js
portal.rpc.serve({handlers: {echo: msg => msg.split(" ").reverse().join(" ")}}).catch(portal.log)
```

This will register echo handler on dedicated **Astral Port** with name prefixed by application package.

Then you can call echo handler from frontend by including the following script
in [index.html](./example/echo-portal/gui/index.html).

```html

<script>
    const service = portal.rpc.bind({"new.portal.service": ["echo"]})
    service.echo("hello astral").then(alert).catch(console.log)
</script>
```

After Saving modification:
1. The application window should reload automatically. 
2. Frontend should send message to the service.
3. Service should respond with reversed message. 
4. Frontend should display response in popup.

NOTE: By default, headless apps close automatically with a 5-second timeout when there are no ongoing connections.
**Portal** can start the closed applications on demand when another app tries to communicate.

If you don't want the app to close automatically, include a following code
in [portal.json](./example/echo-portal/service/portal.json)

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

To build application bundle ready to install on the user environment, use `portal-dev b`, for example:

```shell
portal-dev b ./example/echo-portal
```

This command will search the given directory for projects generating bundles into [build](./example/echo-portal/build) directory.
