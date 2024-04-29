# TODO

* Spawn multiple applications as separate portal processes.
  * [x] prod
  * dev
* [x] Local apps service for testing needed until module isn't available.
  * [x] Features
    * [x] Install app
    * [x] Uninstall app
    * Apps page cursor - allows to load & app changes as well as scroll through the list.
* Invoke queried service on demand. If hosing application is not running spawn it and proceed the query.
  * [x] prod
  * dev
* [x] Close backend with given timeout when idling detected.
* [x] Portal launcher tray icon.
  * [x] Bind portal serve on launch.
* [x] Implement portal installer.
* Upgrade JavaScript RPC:
  * Register methods automatically under the service name specified in manifest.
  * Ensure or redesign. 
* [x] Migrate go-apphost-jrpc into portal repository.
* Add cross integration rpc tests for js - golang.
* Design runtime API for JavaScript apps:
  * backend
  * wails
  * android
* Install JavaScript runtime dependencies in local node modules.
* Use injected prefixed logger. Review & unify existing logs.
* [x] Close nested apps gracefully.