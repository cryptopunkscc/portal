/* eslint-disable */
const _default_bindings = () => ({
  astral_conn_accept: _astral_conn_accept,
  astral_conn_close: _astral_conn_close,
  astral_conn_read: _astral_conn_read,
  astral_conn_write: _astral_conn_write,
  astral_node_info: _astral_node_info,
  astral_query: _astral_query,
  astral_query_name: _astral_query_name,
  astral_resolve: _astral_resolve,
  astral_service_close: _astral_service_close,
  astral_service_register: _astral_service_register,
  sleep: _sleep,
  log: _log,
})

builder.push({
  platform: "default",
  bindings: _default_bindings,
})

const platform = function () {
  for (let next of builder) {
    if (next.platform) {
      return next.platform
    }
  }
}()

const bindings = function () {
  for (let next of builder) {
    if (next.platform) {
      return next.bindings()
    }
  }
}()
