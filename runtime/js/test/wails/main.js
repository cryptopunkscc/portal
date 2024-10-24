import {log, rpc, sleep} from "../../wails.js";
import {runTests} from "../test.js";

runTests("portal.js.test.wails", log, rpc, sleep)
