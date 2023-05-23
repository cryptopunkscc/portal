function astral_conn_accept(arg1) {
  return window['go']['main']['Adapter']['ConnAccept'](arg1);
}

function astral_conn_close(arg1) {
  return window['go']['main']['Adapter']['ConnClose'](arg1);
}

function astral_conn_read(arg1) {
  return window['go']['main']['Adapter']['ConnRead'](arg1);
}

function astral_conn_write(arg1, arg2) {
  return window['go']['main']['Adapter']['ConnWrite'](arg1, arg2);
}

function astral_node_info(arg1) {
  return window['go']['main']['Adapter']['NodeInfo'](arg1);
}

function astral_query(arg1, arg2) {
  return window['go']['main']['Adapter']['Query'](arg1, arg2);
}

function astral_query_name(arg1, arg2) {
  return window['go']['main']['Adapter']['QueryName'](arg1, arg2);
}

function astral_resolve(arg1) {
  return window['go']['main']['Adapter']['Resolve'](arg1);
}

function astral_service_close(arg1) {
  return window['go']['main']['Adapter']['ServiceClose'](arg1);
}

function astral_service_register(arg1) {
  return window['go']['main']['Adapter']['ServiceRegister'](arg1);
}

function log(...arg1) {
  return window['go']['main']['Adapter']['LogArr'](arg1);
}

function sleep(arg1) {
  return window['go']['main']['Adapter']['Sleep'](arg1);
}
