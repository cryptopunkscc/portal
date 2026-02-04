import {writable} from 'svelte/store';

function readableSet(idKey = (value) => value.id, initial = []) {
  const map = new Map(initial.map(value => [idKey(value) ?? Symbol(), value]))
  const store = writable(Array.from(map.values()))

  let set = (...values) => {
    let changed = false
    for (const value of values) {
      let key = idKey(value)
      if (!key) key = Symbol()
      else if (map.get(key) === value) continue
      map.set(key, value)
      changed = true
    }
    if (changed) store.set(Array.from(map.values()))
  }
  set.subscribe = store.subscribe
  set.set = set

  return set
}

export default readableSet