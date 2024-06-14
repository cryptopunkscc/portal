import {bindings} from "./bindings.js";
import {ApphostClient} from "./apphost/adapter.js";
import {RpcClient} from "./rpc/client.js";
import "./apphost/rpc-client.js";
import "./apphost/rpc-service.js";
import "./apphost/rpc-conn.js";
import "./apphost/rpc-handler.js"

export const {log, sleep, platform} = bindings
export const apphost = new RpcClient();
export const rpc = new RpcClient();
