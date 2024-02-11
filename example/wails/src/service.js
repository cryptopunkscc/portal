class Service {

    constructor() {
        this.counter = 0
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

    async inc() {
        return ++this.counter
    }
}

appHost.bindRpc(Service, "rpc").catch(log)
