# Android Astral Agent

[astrald](https://github.com/cryptopunkscc/astrald) agent for Android OS

## Features

* Foreground service for keeping astrald alive
* Ongoing notification
* Embedded goja runner for js app backend
* WebView runner for js app frontend
* Apphost support via js adapter
* Installing js app from zip bundle
* Node log screen
* Node config editor
* Admin panel console

## Dev dependencies

* GO 1.21
* OpenJDK 17
* NDK
* [gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile)

## How to build

Build app & generate apk file for Android OS:

```shell
./gradlew :app:assembleDebug
```

Locate generated apk file under following path:

```shell
./app/build/outputs/apk/debug/app-debug.apk
```

## How to install

Just copy apk to Android device manually and install.

Or install & start using adb commands:

```shell
adb install ./app/build/outputs/apk/debug/app-debug.apk
adb shell am start -n cc.cryptopunks.portal/cc.cryptopunks.portal.MainActivity
```

## How to create JS app

To create a JS app compatible with the android agent, please refer to the
following [manual](https://github.com/cryptopunkscc/js-apphost-adapter/blob/master/example/react.md#create-bundle-with-frontend--backend)
and example [react-basic](https://github.com/cryptopunkscc/js-apphost-adapter/tree/master/example/react-basic) app.
