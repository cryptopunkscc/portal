import {log, rpc} from "portal";

log("launcher start")
const client = rpc.bind({
  "portal": ["open"],
  "cc.cryptopunks.portal.apps": ["install", "uninstall"],
})
client.observe = async () => {
  const conn = await rpc.conn("cc.cryptopunks.apps.observe")
  return {
    next: async () => await conn.decode(),
    more: async (num) => await conn.encode(num),
    close: conn.close,
  }
}

export default client
