export function runTests(port, portal) {
  const log = portal.log
  const rpc = portal.rpc
  const sleep = portal.sleep

  rpc.serve({
    inject: {
      state: {
        counter: 0
      }
    },
    routes: [
      "*",
    ],
    handlers: {
      func0: () => 0,
      func1: ({$, ...opts}, ...args) => [opts, ...args],
      func2: async ({$: {conn, state}}, initial, max) => {
        state.counter = initial
        // call decode from promise to intercept connection close on client side ASAP
        new Promise(() => conn.decode()).finally()
        while (!conn.done && state.counter <= max) {
          await conn.encode(state.counter++)
          await sleep(1)
        }
        return {Type: "eos"}
      },
      func3: (_, msg) => {
        throw msg
      },
      func4: (_, msg) => {
        throw msg
      },
      func5: () => undefined,
      func6: async ({$: {conn}}) => {
        for (let i = 0; i < 10; i++) {
          await sleep(5)
          await conn.write([i])
        }
        await conn.close()
      },
    },
  }).catch(log)

  const methods = [
    "func0",
    "func1",
    "*func2",
    "func3",
    "*func4",
    "func5",
    "func6",
    "func7",
  ]

  async function test() {
    let failed = 0
    await sleep(200)

    // const conn = await rpc.conn(port)
    // const client = conn.bind(...methods)
    const client = rpc.target(port).bind(...methods)

    await log("\n\n\n")
    await log("====================== TEST BEGIN ======================\n")
    await test0()
    await test1()
    await test2()
    await test2_1()
    await test3()
    await test4()
    await test5()
    await test6()
    await log("====================== TEST END ======================\n\n\n")
    await rpc.interrupt()

    async function test0() {
      const expected = 0
      const actual = await client.func0()

      await assert("test0", expected, actual)
    }

    async function test1() {
      const expected = [{a: 1, b: 2}, "asd", true]
      const actual = await client.func1(...expected)

      await assert("test1", expected, actual)
    }

    async function test2() {
      const initial = 3
      const max = 10
      const expected = 5

      await log(`test2 start ${client.func2.request}`)
      const actual = await client.func2.filter(next => next === expected).request(initial, max)
      await assert("test2", expected, actual)
    }

    async function test2_1() {
      const initial = 1
      const max = 3
      const expected = [1, 4, 9]

      await log(`test2_1 start ${client.func2.request}`)
      const actual = await client.func2.map(next => next * next).collect(initial, max)
      await assert("test2_1", expected, actual)
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

    async function test6() {
      const conn = await client.func6.conn()
      let expected = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
      let actual = []
      const buffer = new ArrayBuffer(4)
      const view = new Uint8Array(buffer)
      let len = -1
      do {
        await sleep(Math.floor(Math.random() * 30) + 5)
        len = await conn.read(buffer)
        for (let i = 0; i < len; i++) {
          actual = [...actual, view[i]]
        }
      } while (len > -1)
      await assert("test6", expected, actual)
    }

    async function assert(test, expected, actual) {
      expected = JSON.stringify(expected)
      actual = JSON.stringify(actual)
      let result = `[PASSED] ${test}\n`
      if (expected !== actual) {
        failed++
        result = `[FAILED] ${test}\nactual:\t\t${actual}\nexpected:\t${expected}\n`
      }
      await log(result)
    }

    return failed
  }

  return test().catch(log)
}
