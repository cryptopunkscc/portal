class Service {

    async get(_, arg) {
        return {
            arg: arg,
            val: "Hello RPC"
        }
    }

    async sum(_, a, b) {
        await portal.log(JSON.stringify(a))
        return a + b
    }
}

portal.log("start backend")
portal.rpc.serve({
    handlers: new Service(),
    routes: ["*"],
}).catch(portal.log)