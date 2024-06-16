import {inject} from "../bindings";

// ================== Wails bindings adapter ==================


let wails
try {
  wails = window
} catch (e) {
  wails = {}
}

/* eslint-disable */
const platform = typeof wails['go'] === "undefined" ? undefined : "wails"

/* eslint-disable */
const adapter = () => ({
  // apphost
  astral_conn_accept: wails['go']['main']['Adapter']['ConnAccept'],
  astral_conn_close: wails['go']['main']['Adapter']['ConnClose'],
  astral_conn_read: wails['go']['main']['Adapter']['ConnRead'],
  astral_conn_write: wails['go']['main']['Adapter']['ConnWrite'],
  astral_node_info: wails['go']['main']['Adapter']['NodeInfo'],
  astral_query: wails['go']['main']['Adapter']['Query'],
  astral_query_name: wails['go']['main']['Adapter']['QueryName'],
  astral_resolve: wails['go']['main']['Adapter']['Resolve'],
  astral_service_close: wails['go']['main']['Adapter']['ServiceClose'],
  astral_service_register: wails['go']['main']['Adapter']['ServiceRegister'],
  astral_interrupt: wails['go']['main']['Adapter']['Interrupt'],
  // runtime
  sleep: wails['go']['main']['Adapter']['Sleep'],
  log: async (...arg) => await wails['go']['main']['Adapter']['LogArr'](arg),
})

inject(platform, adapter)
