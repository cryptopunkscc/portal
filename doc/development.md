## App Development

### Dependencies

* [npm](https://nodejs.org/en/download/package-manager) - For npm based templates like `svelte` or `js-rollup`.
* [go](https://go.dev/doc/install) - For `golang` apps.

### Creating

#### HTML app

```shell
portal dev.html new -template html <name|path>
```

This will generate new raw html-based GUI application.

```shell
portal dev.html new -template svelte <name|path>
```

This will generate new svelte-based HTML project.

```shell
portal dev.html templates 
```

This will list supported html-based templates.

#### JS service

```shell
portal dev.js new -template js <name|path>
```

This will generate new raw ES5/ES6-based service.

```shell
portal dev.js new -template js-rollup <name|path>
```

This will generate new rollup-based JS project.

```shell
portal dev.js templates 
```

This will list supported js-compatible templates.

### Developing

Basically, open generated project in your preferred IDE, and code.

Check JSDoc in [api.js](../core/js/src/api.js) to review exported bindings to the native apphost & portal API.   

### Running

```shell
portal -d <path>
```

This will run all apps found recursively in the given directory using the development runner with hot reloading support.

### Building

```shell
portal build -p <path>
```

This will build and pack all apps found recursively in a given directory.

### Publishing

```shell
portal app publish <path>
```

This will publish all app bundles found in given local storage path.
