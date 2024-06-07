// const apphost = portal.apphost
const log = portal.log
const apphost = portal.apphost
const sleep = portal.sleep

const func1 = apphost.rpcQuery("", "test.go.request.func1")
const func2 = apphost.rpcQuery("", "test.go.request.func2")
const func3 = apphost.rpcQuery("", "test.go.request.func3")
const func4 = apphost.rpcQuery("", "test.go.request.func4")

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
  const actual = await func1(expected, false)
  assertEqual(test_func1_a, expected, actual)
}

async function test_func1_b() {
  const message = "message"
  const expected = {error: message}
  const actual = await func1(message, true)

  assertEqual(this, expected, actual)
}

async function test_func2() {
  const expected = [true, 1, 99.99, "text"]
  const actual = await func2(...expected)

  assertEqual(this, expected, actual)
}

async function test_func3_a() {
  const expected = {struct1: {b: true, i: 1, f: 99.99, s: "text"}}
  const actual = await func3(expected)

  assertEqual(this, expected, actual)
}

async function test_func3_b() {
  const expected = null
  const actual = await func3(expected)

  assertEqual(this, expected, actual)
}

async function test_func4() {
  const arg = [true, 1, 99.99, "text"]
  const expected = {b: true, i: 1, f: 99.99, s: "text"}
  const actual = await func4(...arg)

  assertEqual(this, expected, actual)
}

const error = (...args) => portal.log("[FAILED]:", ...args)

function assertEqual(f, l, r) {
  l = JSON.stringify(l)
  r = JSON.stringify(r)
  log("[ASSERT]", f.name, "=====>", l, "equal", r)
  if (l !== r) {
    throw f.name + ": " + l + " != " + r
  }
}

async function main() {
  await sleep(200)
  await Promise.allSettled(tests.map(function (value) {
    return value.call(value).catch(error)
  }))
}

main().catch(log)
