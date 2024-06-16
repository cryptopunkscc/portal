import {rpc, log, sleep} from "../common";

rpc.serve({
  inject: {
    state: {
      counter: 0
    }
  },
  handlers: {
    func0: () => 0,
    func1: (arg) => arg,
    func2: async (initial, max, {conn, state}) => {
      state.counter = initial
      // call decode from promise to intercept connection close on client side ASAP
      new Promise(() => conn.decode()).finally()
      while (!conn.done && state.counter <= max) {
        await conn.encode(state.counter++)
        await sleep(1)
      }
    },
  },
}).catch(log)

test().catch(log)

const client = rpc.bind("portal.js.test", "func0", "func1", "*func2")

async function test() {
  await sleep(200)
  await test0()
  await test1()
  await test2()
  await rpc.interrupt()
}

async function test0() {
  const expected = 0
  const actual = await client.func0()

  await assert("test0", expected, actual)
}

async function test1() {

  const expected = ["a", 1, true, [1, 2, 3], {b: "b"}]
  const actual = await client.func1(expected)

  await assert("test1", expected, actual)
}

async function test2() {
  const initial = 3
  const max = 10
  const expected = 5
  const actual = await client.func2(initial, max, next => {
    return next === expected ? next : null
  })
  await assert("test2", expected, actual)
}

async function assert(test, expected, actual) {
  expected = JSON.stringify(expected)
  actual = JSON.stringify(actual)
  let result = `${test} PASSED`
  if (expected !== actual) {
    result = `${test} actual !== expected:\n${actual}\n${expected}\n`
  }
  await log(result)
}