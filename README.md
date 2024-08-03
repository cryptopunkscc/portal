# Portal to Astral

**Portal** is a runtime & development environment for multiplatform decentralized applications driven
by **Astral**.

### [What is the Astral Network](https://github.com/cryptopunkscc/astrald/blob/master/docs/overview.md)

> Astral is an abstract network that provides authenticated and encrypted connections over a variety of physical
> networks.
> It provides simple and secure connectivity interface, which automatically adapts to existing network conditions.
> Its mission is to dramatically reduce the time it takes to build robust peer-to-peer networks.

### What are the Astral Apps

**Astral Apps** are any applications capable to connect with **Astral** via
[**Apphost Protocol**](https://github.com/cryptopunkscc/astrald/blob/master/mod/apphost/proto/protocol.md).

[**Apphost Library**](https://github.com/cryptopunkscc/astrald/tree/master/lib/astral) is natively written in go, so you
can import it into Golang project and write **Astral App**.
While this approach can be useful for some cases, most likely may not be very convenient.

### Why Portal

**Portal** is aimed to provide runtime environment for multiplatform **Astral Apps** and simplify the development
process.

## How to use

Depending on your case you may want to run the **Astral App**, develop the new one, or compile **Portal Project**.
Read the following docs to learn more about possible use cases.

* [User guide](./user.md) - How to **install** and use **Portal Environment** for users.
* [Developer guide](./developer.md) - How to create, develop, and build **Astral Apps** using **Portal**.
* [Contributor guide](./contributor.md) - How to compile repository and general installer.
* [Technical overview](./doc/overview.md) - Explains some technical details about the project.

# Current Status

This project is at the **alpha** stage, which means is ready for testing and developing proof of concept apps but the
SDK isn't complete and the API may change before the beta release.

## Supported Platforms

List of platforms planned to include in support:

* [x] Linux
    * [x] Debian
    * [ ] Others - Not tested. Should work out of the box, but may require to install some dependencies manually.
      See [tray](cmd/make/tray.go) and [webview](cmd/make/wails.go)
* [x] Windows
* [ ] Android - Outdated PoC, require adjustments
* [ ] macOS - require adjustments in installer
* [ ] iOS - TODO
