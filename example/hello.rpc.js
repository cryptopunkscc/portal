class Service {

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

appHost.bindRpc(Service, "srv").catch(log)
