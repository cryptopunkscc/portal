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

portal.log("start backend")
portal.rpc.serve({
    handlers: new Service(),
    routes: ["*"],
}).catch(portal.log)