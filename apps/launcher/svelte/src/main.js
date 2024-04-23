import './style.css'
import Launcher from "./Launcher.svelte";

const app = new Launcher({
  target: document.getElementById('app')
})

export default app
