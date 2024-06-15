import {inject} from "../bindings";

// ================== Wails bindings adapter ==================

/* eslint-disable */
const platform = typeof window['go'] === "undefined" ? undefined : "wails"

/* eslint-disable */
const adapter = () => ({
  // apphost
  astral_conn_accept: window['go']['main']['Adapter']['ConnAccept'],
  astral_conn_close: window['go']['main']['Adapter']['ConnClose'],
  astral_conn_read: window['go']['main']['Adapter']['ConnRead'],
  astral_conn_write: window['go']['main']['Adapter']['ConnWrite'],
  astral_node_info: window['go']['main']['Adapter']['NodeInfo'],
  astral_query: window['go']['main']['Adapter']['Query'],
  astral_query_name: window['go']['main']['Adapter']['QueryName'],
  astral_resolve: window['go']['main']['Adapter']['Resolve'],
  astral_service_close: window['go']['main']['Adapter']['ServiceClose'],
  astral_service_register: window['go']['main']['Adapter']['ServiceRegister'],
  astral_interrupt: window['go']['main']['Adapter']['Interrupt'],
  // runtime
  sleep: window['go']['main']['Adapter']['Sleep'],
  log: async (...arg) => await window['go']['main']['Adapter']['LogArr'](arg),
})

inject(platform, adapter)
