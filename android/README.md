# Android Astral Agent

[astrald](https://github.com/cryptopunkscc/astrald) agent for Android OS

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

## Android Studio

Compiling the project via Android Studio might not work by default because of different env and $PATH.
This problem might be solved by running AS from the terminal.
