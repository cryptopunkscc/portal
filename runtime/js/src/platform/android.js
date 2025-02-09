import {inject} from "../bindings";
// ================== Android bindings adapter ==================

/* eslint-disable */
const platform = typeof _app_host === "undefined" ? undefined : "android"

/* eslint-disable */
const adapter = () => {

  const _awaiting = new Map()

  window._resolve = (id, value) => {
    _awaiting.get(id)[0](value)
    _awaiting.delete(id)
  }

  window._reject = (id, error) => {
    _awaiting.get(id)[1](error)
    _awaiting.delete(id)
  }

  const _promise = (block) =>
    new Promise((resolve, reject) =>
      _awaiting.set(block(), [resolve, reject]))

  return {
    // apphost
    astral_node_info: (arg1) => _promise(() => _app_host.nodeInfo(arg1)).then(v => JSON.parse(v)),
    astral_conn_accept: (arg1) => _promise(() => _app_host.connAccept(arg1)),
    astral_conn_close: (arg1) => _promise(() => _app_host.connClose(arg1)),
    astral_conn_read: (arg1, arg2) => _promise(() => {
      // TODO write result to byte array
      return _app_host.connRead(arg1, arg2);
    }),
    astral_conn_write: (arg1, arg2) => _promise(() => _app_host.connWrite(arg1, arg2)),
    astral_conn_read_ln: (arg1) => _promise(() => _app_host.connReadLn(arg1)),
    astral_conn_write_ln: (arg1, arg2) => _promise(() => _app_host.connWriteLn(arg1, arg2)),
    astral_query: (arg1, arg2) => _promise(() => _app_host.query(arg1, arg2)),
    astral_query_name: (arg1, arg2) => _promise(() => _app_host.queryName(arg1, arg2)),
    astral_resolve: (arg1) => _promise(() => _app_host.resolve(arg1)),
    astral_service_close: (arg1) => _promise(() => _app_host.serviceClose(arg1)),
    astral_service_register: (arg1) => _promise(() => _app_host.serviceRegister(arg1)),
    astral_interrupt: () => _promise(() => _app_host.interrupt()),
    // runtime
    sleep: (arg1) => _promise(() => _app_host.sleep(arg1)),
    log: (arg1) => _app_host.log(arg1),
    exit: (arg1) => _app_host.exit(arg1),
  }
}

inject(platform, adapter)
