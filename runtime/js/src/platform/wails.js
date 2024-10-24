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
  astral_conn_read: async (id, buffer) => {
    const base64 = await wails['go']['main']['Adapter']['ConnRead'](id, buffer.byteLength);
    const binary = atob(base64);
    const len = binary.length;
    const view = new Uint8Array(buffer)
    for (let i = 0; i < len; i++) {
      view[i] = binary.charCodeAt(i)
    }
    return len
  },
  astral_conn_write: wails['go']['main']['Adapter']['ConnWrite'],
  astral_conn_read_ln: wails['go']['main']['Adapter']['ConnReadLn'],
  astral_conn_write_ln: wails['go']['main']['Adapter']['ConnWriteLn'],
  astral_node_info: wails['go']['main']['Adapter']['NodeInfo'],
  astral_query: wails['go']['main']['Adapter']['Query'],
  astral_query_name: wails['go']['main']['Adapter']['QueryName'],
  astral_resolve: wails['go']['main']['Adapter']['Resolve'],
  astral_service_close: wails['go']['main']['Adapter']['ServiceClose'],
  astral_service_register: wails['go']['main']['Adapter']['ServiceRegister'],
  astral_interrupt: wails['go']['main']['Adapter']['Interrupt'],
  // runtime
  sleep: wails['go']['main']['Adapter']['Sleep'],
  log: wails['go']['main']['Adapter']['Log'],
})

inject(platform, adapter)
