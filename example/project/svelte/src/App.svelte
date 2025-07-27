<script>
  import {apphost, log, platform, rpc} from 'portal';
  import {onDestroy, onMount} from "svelte";

  const service = rpc.target("example.project.js").bind("link", "*ticker")

  const max = 10

  const state = {
    info: "undefined",
    status: "started",
    limit: 3,
    delay: 1000,
  }

  onMount(async () => {
    connect().catch(e => {
      state.status = e
      console.log(e)
    })
    info().catch(log)

  })

  onDestroy(() => {
    rpc.interrupt()
  })


  async function connect() {
    state.status = "connecting..."
    await service.link.conn()
    state.status = "connected"
  }

  /**
   * Example of registering selector.
   */
  async function ticker() {
    log(state.limit)
    try {
      state.status = await service.ticker.map(async next => {
        state.status = next
        if (next >= max) return "max"
      }).request(state.limit, state.delay)
      state.status = state.status
    } catch (e) {
      state.status = e
      console.log(e)
      // alert(e)
    }
  }

  /**
   * Example function for resolving node info by name.
   */
  async function info() {
    const id = await apphost.resolve("localnode")
    const info = await apphost.nodeInfo(id)
    state.info = JSON.stringify(info, null, 2)
  }

</script>

<main>
    <p>Portal svelte app - Running on {platform}</p>
    <p>Node info: {state.info}</p>
    <p>Status: {state.status}</p>
    <p><input bind:value={state.limit}/> ticks limit</p>
    <p><input bind:value={state.delay}/> ticks delay</p>
    <div/>
    <button on:click={ticker}>start ticker</button>
</main>


<style>
    body {
        color: white;
    }
</style>