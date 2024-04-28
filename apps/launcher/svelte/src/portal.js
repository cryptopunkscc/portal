import {apphost} from "portal/portal.js";

export default {
  launch: apphost.rpcQuery("", "portal.open"),
  install: apphost.rpcQuery("", "portal.install"),
  uninstall: apphost.rpcQuery("", "postal.uninstall"),
  observe: async (num) => {
    const conn = await apphost.query("", "portal.observe")
    return {
      next: async () => await conn.readJson("observe"),
      more: async (num) => await conn.writeJson(num),
      close: conn.close,
    }
  }
}