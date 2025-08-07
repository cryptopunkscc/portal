# Overview

A runtime and management environment for decentralized multiplatform apps.
It provides authentication, identity, and connectivity via [Astral](https://github.com/cryptopunkscc/astrald/blob/master/README.md).
Its mission is to lay the groundwork for building a decentralized network of multiplatform applications.

## Core Design

In basics, it is a tiny daemon that provides compatibility layer between user, device, applications, and Astral.

## Features

* Creates new user.
* Assigns devices to the user.
* Manages applications (installs, authenticates, runs, lists).
* Provides application development tools (building, bundling, publishing).
* Exposes management API.

## Components

Complete portal environment consists of the following core components:

* **[astrald](https://github.com/cryptopunkscc/astrald/tree/master/cmd/astrald)** - A daemon that provides all the networking features like: identity, authentication, connectivity, encryption, and many more.
* **[portald](../cmd/portald)** - A daemon that implements core features for app management. It configures and starts astrald process, and serves management API for the CLI or GUI clients.
* **[portal](../cmd/portal)** - A reference CLI client for portald management API.
* **[installer](../cmd/install-portal-to-astral)** - A platform specific bundle that unpacks mentioned components, installs embedded apps, and initiates user identity.

## Embedded Applications

In addition to the core components, portal delegates the rest of its features to the embedded apps and services:

* **[portal.js](../apps/js)** - A runner for ES5/6 services. Multiplatform, go.
* **[portal.html](../apps/html)** - A runner for HTML5 apps. Implementation may vary depending on platform.
* **[portal.dev.js](../apps/js)** - A runner for creating and developing ES5/6 services. Only for desktop. Requires npm installed on host machine.
* **[portal.dev.html](../apps/html-dev)** - A runner for creating and developing HTML5 apps. Only for desktop. Requires npm installed on host machine.
* **[portal.dev.go](../apps/go-dev)** - A runner for creating and developing go apps. Only for desktop. Requires go installed on host machine.
* **[portal.dev.exec](../apps/exec-dev)** - A proxy runner for developing executable apps. Only for desktop.
* **[portal.launcher](../apps/launcher)** - A launcher for application. Multiplatform, HTML5.
* **[portal.tray](../apps/tray)** - Portal tray icon. Only for desktop.