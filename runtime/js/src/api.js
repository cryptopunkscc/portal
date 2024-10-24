import {bindings} from "./bindings.js";
import {ApphostClient} from "./apphost/adapter.js";
import {RpcClient} from "./rpc/client.js";

export const log = async any => await bindings.log(typeof any == 'object' ? JSON.stringify(any) : any)
export const {sleep, platform} = bindings
export const apphost = new ApphostClient();
export const rpc = new RpcClient();
