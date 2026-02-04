import * as portal from "../../common";
import {runTests} from "../test";

runTests("self", portal).then(code => portal.exit(code))
