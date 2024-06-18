import {log, sleep} from "portal";

export default {
  launch: async (id) => log("PortalMock.launch: " + id),
  install: async (id) => log("PortalMock.install: " + id),
  uninstall: async (id) => log("PortalMock.uninstall: " + id),
  observe: async () => {
    let buff = []
    let end = 0
    return {
      next: async () => {
        for (; ;) {
          if (buff.length === 0) {
            await sleep(100)
            continue
          }
          // log(buff)
          return buff.splice(0, 1)[0]
        }
      },
      more: async (num) => {
        // log("more")
        const begin = end
        end = begin + num
        // log(begin, end)
        for (let i = begin; i < end; i++) {
          let next = {id: i, name: "app" + i}
          buff.push(next)
          // log("next: ", next)
        }
      }
    }
  }
}