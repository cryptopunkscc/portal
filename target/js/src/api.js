import {bindings} from "./bindings.js";
import {ApphostClient} from "./apphost/adapter.js";
import "./apphost/rpc.js";

export const {log, sleep, platform} = bindings
export const apphost = new ApphostClient();
