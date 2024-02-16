import {useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import appHost, {log, platform} from './apphost';

function App() {
  log("render")

  const [rpc, setRpc] = useState({})
  const [state, setState] = useState({
    info: "undefined",
    sum: 0,
    counter: 0,
  })

  useEffect(() => {
    async function connect() {
      try {
        log("rpc connecting...")
        const conn = await appHost.bindRpc("", "rpc")
        // let conn = await appHost.query("", "rpc")
        // await conn.bindRpc()
        setRpc(conn)
        log("rpc connected")
      } catch (e) {
        log(e)
      }
    }
    connect().catch(log)
  }, [])

  async function info() {
    const id = await appHost.resolve("localnode")
    const info = await appHost.nodeInfo(id)
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
