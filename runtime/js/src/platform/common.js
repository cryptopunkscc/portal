import {inject} from "../bindings";

const platform = typeof _log === 'undefined' ? undefined : "common"

/* eslint-disable */
const adapter = () => ({
  // apphost
  astral_conn_accept: _astral_conn_accept,
  astral_conn_close: _astral_conn_close,
  astral_conn_read: async (id, buffer) => {
    const array = await _astral_conn_read(id, buffer.byteLength);
    const len = array.length;
    const view = new Uint8Array(buffer);
    for (let i = 0; i < len; i++) {
      view[i] = array[i];
    }
    return len;
  },
  astral_conn_write: _astral_conn_write,
  astral_conn_read_ln: _astral_conn_read_ln,
  astral_conn_write_ln: _astral_conn_write_ln,
  astral_node_info: _astral_node_info,
  astral_query: _astral_query,
  astral_query_name: _astral_query_name,
  astral_resolve: _astral_resolve,
  astral_service_close: _astral_service_close,
  astral_service_register: _astral_service_register,
  astral_interrupt: _astral_interrupt,
  // apphost
  sleep: _sleep,
  log: _log,
})

inject(platform, adapter)
