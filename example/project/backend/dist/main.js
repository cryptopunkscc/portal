/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
var __webpack_exports__ = {};

;// CONCATENATED MODULE: ./src/service.js
class Service {
  constructor() {
    this.name = "rpc";
    this.counter = 0;
  }
  async get(arg) {
    return {
      arg: arg,
      val: "Hello RPC"
    };
  }
  async sum(a, b) {
    return a + b;
  }
  async inc() {
    return ++this.counter;
  }
  async ticker(send) {
    let counter = 0;
    for (;;) {
      await send(counter++);
      await sleep(1000);
    }
  }
}
;// CONCATENATED MODULE: ./index.js
// import apphost from "../../lib/apphost/apphost";

appHost.bindRpcService(Service);
/******/ })()
;