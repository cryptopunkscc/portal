const {apphost, log, sleep} = portal

log("start backend")
class Service {

    constructor() {
        this.name = ""
    }

    async get(arg) {
        return {
            arg: arg,
            val: "Hello RPC"
        }
    }

    async sum(a, b) {
        return a + b
    }
}

async function bind() {
    apphost.bindRpcService(Service).catch(log)
}

bind().catch(log)
