import {bindings} from "./bindings.js";
import {ApphostClient} from "./apphost/adapter.js";
import {RpcClient} from "./rpc/client.js";

export const {log, sleep, platform} = bindings
export const apphost = new ApphostClient();
export const rpc = new RpcClient();
