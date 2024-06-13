# TODO

* Spawn multiple applications as separate portal processes.
    * [x] prod
    * [x] dev
* [x] Local apps service for testing needed until module isn't available.
    * [x] Features
        * [x] Install app
        * [x] Uninstall app
        * Apps page cursor - allows to load & app changes as well as scroll through the list.
* [x] Invoke queried service on demand. If hosing application is not running spawn it and proceed the query.
    * [x] prod
    * [x] dev
* [x] Close backend with given timeout when idling detected.
* [x] Portal launcher tray icon.
    * [x] Bind portal serve on launch.
* [x] Implement portal installer.
* [x] Upgrade JavaScript RPC:
    * [x] Register methods automatically under the service name specified in manifest.
    * [x] Ensure or redesign.
* [x] Migrate go-apphost-jrpc into portal repository.
* [x] Install JavaScript runtime dependencies in local node modules.
* [x] Use injected prefixed logger. Review & unify existing logs.
* [x] Close nested apps gracefully.
* [x] Cross integration rpc tests for js - golang.
* [x] Golang runner support.
* Design runtime API for JavaScript apps:
    * backend
    * wails
    * android

## Bugs

* Fix backend timeout.
* Fix go runners reload on change.