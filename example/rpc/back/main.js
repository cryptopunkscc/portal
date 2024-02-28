class Service {

    constructor() {
        this.name = "srv"
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

appHost.bindRpcService(Service).catch(log)
