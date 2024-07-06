// const apphost = portal.apphost
const log = portal.log
const rpc = portal.rpc
const sleep = portal.sleep

const services = [
  "go",
  "js",
]

const methods = [
  "func1",
  "func2",
  "func3",
  "func4",
]

const flows = async service => {
  const conn = await rpc.conn(`test.${service}.flow`)
  return conn.bind(...methods)
}

const requests = async service => {
  const api = {}
  api[`test.${service}.request`] = methods
  return await rpc.bind(api);
}

const connections = [
  flows,
  requests,
]

const tests = [
  test_func1_a,
  test_func1_b,
  test_func2,
  test_func3_a,
  test_func3_b,
  test_func4,
]

async function test_func1_a() {
  const expected = "message"
  const actual = await this.func1(expected, false)
  assertEqual(this, expected, actual)
}

async function test_func1_b() {
  const expected = "message"
  let actual
  try {
    await this.func1(expected, true)
  } catch (e) {
    actual = e
  }
  assertEqual(this, expected, actual)
}

async function test_func2() {
  const expected = [true, 1, 99.99, "text"]
  const actual = await this.func2(...expected)
  assertEqual(this, expected, actual)
}

async function test_func3_a() {
  const expected = {struct1: {b: true, i: 1, f: 99.99, s: "text"}}
  const actual = await this.func3(expected)
  assertEqual(this, expected, actual)
}

async function test_func3_b() {
  const expected = null
  const actual = await this.func3(expected)
  assertEqual(this, expected, actual)
}

async function test_func4() {
  const arg = [true, 1, 99.99, "text"]
  const expected = {b: true, i: 1, f: 99.99, s: "text"}
  const actual = await this.func4(...arg)
  assertEqual(this, expected, actual)
}

const error = (e) => portal.log(`FAILED ${e}`)

function assertEqual(f, l, r) {
  l = JSON.stringify(l)
  r = JSON.stringify(r)
  if (l !== r) {
    throw f.name + " " + l + " != " + r
  } else {
    // log(`TEST_ASSERT`, `${f.name}(${l}==${r})`)
  }
}

async function main() {
  await sleep(200)
  for (let service of services) {
    for (let connection of connections) {
      const conn = await connection(service)
      for (let test of tests) {
        try {
          await test.call(Object.assign(test, conn))
          log(`PASSED ${test.name}`)
        } catch (e) {
          error(e)
        }
      }
    }
  }
}

main().catch(error)
