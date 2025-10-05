import {rpc, log} from 'portal';
import {writable} from 'svelte/store';

const method = "user.check"
const key = method
function hasUser() {
  const stored = localStorage.getItem(key)
  const data = stored ? JSON.parse(stored) : null
  const store = writable(data);
  let refreshing = false

  store.subscribe(value => {
    if (value) {
      const string = JSON.stringify(value)
      localStorage.setItem(key, string);
    } else {
      localStorage.removeItem(key);
    }
  })

  store.refresh = async () => {
    if (refreshing) return
    refreshing = true
    try {
      const has = await rpc.target("portald").call(key).request()
      store.set(has)
    } catch (e) {
      log(e)
      store.set(undefined)
    } finally {
      refreshing = false
    }
  }

  store.refresh()

  return store;
}

export default hasUser()
