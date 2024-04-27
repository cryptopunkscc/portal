# TODO

* Spawn multiple applications as separate portal processes.
  * [x] prod
  * dev
* [x] Local apps service for testing needed until module isn't available.
  * Features
    * Install app
    * Uninstall app
    * Apps page cursor - allows to load & app changes as well as scroll through the list.
* Invoke queried service on demand. If hosing application is not running spawn it and proceed the query.
  * [x] prod
  * dev
* Close backend with given timeout when idling detected.
* Implement a cmd for building portal app along with js scripts.
* Upgrade JavaScript RPC:
  * Register methods automatically under the service name specified in manifest.
  * Ensure or redesign. 
* Migrate go-apphost-jrpc into portal repository.
* Add cross integration rpc tests for js - golang.
* Design runtime API for JavaScript apps:
  * backend
  * wails
  * android
