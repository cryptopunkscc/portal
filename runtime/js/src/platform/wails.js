import {inject} from "../bindings";

// ================== Wails bindings adapter ==================


let wails
try {
  wails = window['go']['main']['Adapter']
} catch {
}

/* eslint-disable */
const platform = wails ? "wails" : undefined

/* eslint-disable */
const adapter = () => ({
  // apphost
  astral_conn_accept: wails['ConnAccept'],
  astral_conn_close: wails['ConnClose'],
  astral_conn_read: async (id, buffer) => {
    const base64 = await wails['ConnRead'](id, buffer.byteLength);
    const binary = atob(base64);
    const len = binary.length;
    const view = new Uint8Array(buffer)
    for (let i = 0; i < len; i++) {
      view[i] = binary.charCodeAt(i)
    }
    return len
  },
  astral_conn_write: wails['ConnWrite'],
  astral_conn_read_ln: wails['ConnReadLn'],
  astral_conn_write_ln: wails['ConnWriteLn'],
  astral_node_info: wails['NodeInfo'],
  astral_query: wails['Query'],
  astral_resolve: wails['Resolve'],
  astral_service_close: wails['ServiceClose'],
  astral_service_register: wails['ServiceRegister'],
  astral_interrupt: wails['Interrupt'],
  // runtime
  sleep: wails['Sleep'],
  log: wails['Log'],
  exit: wails['Exit'],
})

inject(platform, adapter)
