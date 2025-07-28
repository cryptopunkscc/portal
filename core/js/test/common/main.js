import * as portal from "../../common";
import {runTests} from "../test";

runTests("portal.js.test.common", portal).then(code => portal.exit(code))
