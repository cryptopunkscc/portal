<script>
  import logo from './assets/images/logo-universal.png'
  import apphost, {log, platform} from "../../lib/apphost/apphost";

  init()

  let rpc

  let state = {
    info: "undefined",
    sum: 0,
    counter: 0,
    tick: 0,
  }

  async function init() {
    await apphost.interrupt()
    rpc = await apphost.bindRpc("", "rpc")
    log("rpc connected")
    await subscribeTicker()
  }

  async function subscribeTicker() {
    const read = await rpc.subscribe("ticker")
    for (let i = 0; i < 100; i++) {
      state.tick = await read()
    }
    await read.cancel()
  }

  async function info() {
    const id = await apphost.resolve("localnode")
    const info = await apphost.nodeInfo(id)
    state.info = JSON.stringify(info, null, 2)
    log(state.info)
  }

  async function sum() {
    state.sum = await rpc.sum(2, 2)
  }

  async function inc() {
    state.counter = await rpc.inc()
  }

</script>

<main>
  <img alt="Wails logo" id="logo" src="{logo}">
  <p>Running on {platform}</p>
  <p>{JSON.stringify(rpc)}</p>
  <p>ticker {state.tick}</p>

  <button on:click={info}>get node info</button>
  <p>{state.info}</p>

  <button on:click={sum}>rpc sum 2 + 2</button>
  <p>{state.sum}</p>

  <button on:click={inc}>rpc increment</button>
  <p>{state.counter}</p>
</main>

<style>
  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }
</style>
