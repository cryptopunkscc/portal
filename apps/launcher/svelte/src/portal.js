import {log, rpc} from "portal";

log("launcher start")

const client = rpc.target("portald").bind("open", "install", "uninstall")

client.observe = async () => {
  const conn = await rpc.target("portald").conn("list.observe")
  return {
    next: async () => await conn.decode(),
    more: async (num) => await conn.encode(num),
    close: conn.close,
  }
}

export default client
