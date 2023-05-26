# React quick start

## Install NodeJS

Use your system's package manager or [fnm](https://github.com/Schniz/fnm):

```shell
$ curl -fsSL https://fnm.vercel.app/install | bash
$ fnm install v20.2.0
$ fnm use v20.2.0
```

## Create a new app

Create a new React app. Note that while using the npm development server, the app will not have access to
astral APIs.

```shell
$ npx create-react-app myapp
$ cd myapp
$ npm start
```

## Run using wails runtime

Build the app and launch it using wails runtime:

```shell
$ npm run build
$ astral-runtime-wails build/
```

You can zip the contents of the build directory:
```shell
$ (cd build; zip ../app.zip -r .)
$ astral-runtime-wails app.zip
```