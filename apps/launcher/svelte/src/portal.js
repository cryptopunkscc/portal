import {log, rpc} from "portal";

log("launcher start")
const client = rpc.bind({"portal": ["open", "install", "uninstall"]})
client.observe = async () => {
  const conn = await rpc.conn("portal.observe")
  return {
    next: async () => await conn.decode(),
    more: async (num) => await conn.encode(num),
    close: conn.close,
  }
}

export default client
