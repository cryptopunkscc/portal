import {log, rpc} from "portal/portal";

log("launcher start")
// log(JSON.stringify(rpc), rpc)
const client = rpc.bind("portal", "open", "install", "uninstall")
client.observe = async () => {
  await log("launcher observe")
  const conn = await rpc.query("portal.observe", log)
  await log("launcher observe2")
  return {
    next: async () => await conn.readJson("observe"),
    more: async (num) => await conn.writeJson(num),
    close: conn.close,
  }
}

export default client
