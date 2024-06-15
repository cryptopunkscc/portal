import {rpc, log, sleep} from "../common";

rpc.serve({
  handlers: {
    func1: (arg) => arg,
  }
}).catch(log)

sleep(200)

test().catch(log)

async function test() {
  const client = rpc.bind("portal.js.test", "func1")

  let result = "PASSED"
  let expected = ["a", 1, true, [1, 2, 3], {b: "b"}]
  let actual = await client.func1(expected)

  expected = JSON.stringify(expected)
  actual = JSON.stringify(actual)

  if (expected !== actual) {
    result = `actual != expected:\n${actual}\n${expected}\n`
  }

  log(result)
}