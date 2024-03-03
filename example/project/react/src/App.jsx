import {useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {apphost, log, platform} from '../../lib/portal';

function App() {
  log("render")

  const [rpc, setRpc] = useState({})
  const [tick, setTick] = useState(-1)
  const [state, setState] = useState({
    info: "undefined",
    sum: 0,
    counter: 0,
    tick: 0,
  })

  useEffect(() => {
    async function subscribeTicker(rpc) {
      const read = await rpc.subscribe("ticker")
      for (let i = 0; i < 100; i++) {
        const num = await read()
        setTick(num)
      }
      await read.cancel()
    }

    async function connect() {
      try {
        log("rpc connecting...")
        // const apphost = new AppHostClient()
        const conn = await apphost.bindRpc("", "rpc")
        // const conn = await (new AppHostClient()).bindRpc("", "rpc")
        // let conn = await apphost.query("", "rpc")
        // await conn.bindRpc()
        setRpc(conn)
        log("rpc connected")
        await subscribeTicker(conn)
      } catch (e) {
        log(e)
      }
    }
    connect().catch(log)

    return apphost.interrupt
  }, [])

  async function info() {
    const id = await apphost.resolve("localnode")
    const info = await apphost.nodeInfo(id)
    const string = JSON.stringify(info, null, 2)
    log(string)
    setState({...state, info: string})
  }

  async function sum() {
    const num = await rpc.sum(2, 2)
    setState({...state, sum: num})
  }

  async function inc() {
    const num = await rpc.inc()
    setState({...state, counter: num})
  }


  return (
    <div id="App">
      <img src={logo} id="logo" alt="logo"/>
      <p>Running on {platform}</p>
      <p>{JSON.stringify(rpc)}</p>
      <p>ticker {tick}</p>

      <button onClick={info}>get node info</button>
      <p>{state.info}</p>

      <button onClick={sum}>rpc sum 2 + 2</button>
      <p>{state.sum}</p>

      <button onClick={inc}>rpc increment</button>
      <p>{state.counter}</p>
    </div>
  )
}

export default App
