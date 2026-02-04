import {RpcConn} from "./conn.js";
import {bindings} from "../bindings.js";

// TODO add support async iterators for es5 (goja), most likely using rollup + babel.

const encode = RpcConn.prototype.encode

RpcConn.prototype.encode = async function (data) {
  bindings.log("encode async gen")
  if (isAsyncGenerator(data)) {
    for await (let next of data) {
      await encode.call(this, next)
    }
    return
  }
  await encode.call(this, data)
}

function isAsyncGenerator(any) {
  return "function" === typeof any.next
    && "function" === typeof any[Symbol.asyncIterator]
    && any === any[Symbol.asyncIterator]()
}


RpcConn.prototype.next = async function () {
  if (!this.done) {
    try {
      this.value = await this.decode()
    } catch (e) {
      this.error = e
      this.done = true
    }
  }
  return this
}

RpcConn.prototype.return = async function (value) {
  if (!this.done) {
    if (value) {
      this.value = value
    }
    this.done = true
    this.close().catch()
  }
  return this
}

RpcConn.prototype.throw = async function () {
  if (!this.done) {
    this.done = true
    this.close().catch()
  }
  return this
}
