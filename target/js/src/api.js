import {bindings} from "./bindings.js";
import {ApphostClient} from "./apphost/adapter.js";
import "./apphost/rpc-client.js";
import "./apphost/rpc-service.js";
import "./apphost/rpc-conn.js";

export const {log, sleep, platform} = bindings
export const apphost = new ApphostClient();
