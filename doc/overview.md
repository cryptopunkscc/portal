# How it works

Basically, **Portal** works as a tiny service, capable for running **Astral Apps** as subprocesses and managing them.
In addition, Portal provides a bunch of base general-purpose components.
Each component is a standalone astral app.

### Astral daemon

Portal requires `astrald` to be running on the user space.
For convenience `portal-installer` includes compatible version of `astrald` and installs it along with other dependencies.

Portal itself is capable to start `astrald` as subprocess if needed.
In general, you can start `astrald` as separate process before starting `portal`,
or let `portal` to start `astrald` as a subprocess.

### Application runner

Is a special type of executable component, capable of executing applications written in a dynamically interpreted language on integrated VM.

Our goal is to provide runtime and development environment capable to execute application written in popular language on any possible platform.
Providing a support for dynamically interpreted language seems to be the best way to achieve this goal.

By default, Portal provide first class support for HTML/JS based apps, with limited support for running native executables and developing golang apps.

List of supported languages and available runners may change in the future.



## What is included

Complete portal environment consists of the following executable components:

* `astrald` - A default implementation of Astral network that runs connectivity node on the local machine providing apphost communication interface on tcp or unix sockets. It's a core dependency for portal.
* `portal-installer` - A bundle containing required executables and capable to install them in the user's environment.
* `portal` - A default commandline interface for starting and communicating with `portal-app`.
* `portal-app` - A core service responsible for managing application runners.
* `portal-app-wails` - A HTML webkit runner for desktops driven by wails project.
* `portal-app-goja` - A JS runner driven by `goja` - ES 5.1(+) implementation written in pure go.
* `portal-tray` - Displays tray indicator.
* `portal-dev` - A core service for generating projects, managing development runners, and creating app bundles.
* `portal-dev-wails` - Hot-reloading runner for developing HTML apps driven by wails.
* `portal-dev-goja` - Hot-reloading runner for js apps driven by goja.
* `portal-dev-go` - Hot-reloading runner for golang projects. Depends on `portal-dev-exec`
* `portal-dev-exec` - Hot-reloading runner for executables.
* `anc` - Tool inspired by `netcat`/`nc`. Allows to access apphost interface through command line

Depending on your usecase, you may not need all of them. For example.
* If you are not interested in developing apps you can skip all components prefixed by `portal-dev`.
* If you are running Portal on a headless environment you don't need following UI components:
    * `portal-app-wails`
    * `portal-dev-wails`
    * `portal-tray`