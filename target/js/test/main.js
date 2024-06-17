import {log, rpc, sleep} from "../common";

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
    func3: (msg) => {
      throw msg
    },
    func4: (msg) => {
      throw msg
    },
    func5: () => undefined
  },
}).catch(log)

test().catch(log)

const client = rpc.bind("portal.js.test",
  "func0",
  "func1",
  "*func2",
  "func3",
  "*func4",
  "func5",
)

async function test() {
  await sleep(200)
  log("\n\n\n")
  log("====================== TEST BEGIN ======================\n")
  await test0()
  await test1()
  await test2()
  await test3()
  await test4()
  await test5()
  log("====================== TEST END ======================\n\n\n")
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

async function test3() {
  const expected = "test error"
  let actual
  try {
    actual = await client.func3(expected)
  } catch (e) {
    actual = e
  }
  await assert("test3", expected, actual)
}

async function test4() {
  const expected = "test error"
  let actual
  try {
    actual = await client.func4(expected, () => undefined)
  } catch (e) {
    actual = e
  }
  await assert("test4", expected, actual)
}

async function test5() {
  const expected = {}
  const actual = await client.func5()
  await assert("test5", expected, actual)
}

async function assert(test, expected, actual) {
  expected = JSON.stringify(expected)
  actual = JSON.stringify(actual)
  let result = `[PASSED] ${test}\n`
  if (expected !== actual) {
    result = `[FAILED] ${test}\nactual:\t\t${actual}\nexpected:\t${expected}\n`
  }
  await log(result)
}