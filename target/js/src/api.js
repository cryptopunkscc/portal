import {bindings} from "./bindings";
import {AppHostClient} from "./apphost/client";
import "./apphost/jrpc";

export const {log, sleep, platform} = bindings
export const apphost = new AppHostClient();
