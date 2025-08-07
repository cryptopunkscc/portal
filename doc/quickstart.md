# Quick start

1. [Installation](#Installation)
   1. [First installation](#first-installation)
   2. [Next installation](#next-installation)
2. [App Management](#app-management)
   1. [Running](#running)
   2. [Listing](#listing)
   3. [Installing](#installing)

## Installation

Installation process will copy the following components to local drive:
* astrald - Astral daemon binary
* portald - Portal daemon binary
* portal - Portal CLI binary
* default applications

### Installer binary

First, obtain `install-portal-to-astral` binary. You can either download it or build it from sources:

```shell
git clone https://github.com/cryptopunkscc/portal.git
cd portal
./mage build:installer
```

This will generate `install-portal-to-astral` binary into project's `bin` directory.

#### Linux

Make sure you have `$HOME/.local/bin` added to your `$PATH`:

### First installation

```shell
./install-portal-to-astal <user name>
```

This will install the complete environment and generate a new user identity under the given name. 

### Next installation

You can add more of your devices to a private group.

```shell
./install-portal-to-astal
```

This will install complete environment ready to [claim](#claim-next-device) by existing user.

#### Claim next device

If you already performed [first](#first-installation) and [next](#first-installation) installation on your devices, you can claim your next device from the first one. 

First ensure that both devices are in the same local network, then run following command from the first device:

```shell
portal user claim <node alias|node identity>
```

This will assign the next device to the existing user.

## Running

Close `astrald` process if it is already running, to allow `portald` running own `astrald` process with a custom root path.

Then run: 
```shell
portald
```

This will run portal daemon and astral daemon.

## App Management

### Running

App can be called by its name, package name, or path to containing directory.

```shell
portal <name|package|path>
```

Portal provides few preinstalled apps like launcher or tray indicator. 

For example:

```shell
portal launcher
```

will run portal apps launcher by its name.

```shell
portal portal.tray
```

will run portal tray icon by its package.

```shell
portal ./example/rpc
```

will run example client and service stored in the ./example/rpc directory.

### Listing

```shell
portal apps list
```

This will list installed apps. 

```shell
portal apps available
```

This will list apps available to install via network. *Currently, only users devices group is supported*

```shell
portal help
```

The help command also prints installed apps.

### Installing

```shell
portal app install <name|package|bundleId|path>
```

This will install application from network or local storage depending on provided argument. 
If the argument is a local storage path, this command recursively searches given path installing all applications found.

### What's next?

* [App Development](./development.md)
