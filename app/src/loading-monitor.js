import { writable } from 'svelte/store'

const loader = writable({
  m: new Map(),
  i: 0,
  dirty: false,
})

let finishedLoading = false

export const loading = (id = undefined, count = 1) => {
  if (finishedLoading || !process.browser) {
    return () => {}
  }
  let key = undefined
  loader.update(l => {
    key = id === undefined ? l.i++ : `${id}-${l.i++}`
    l.m.set(key, count)
    if (!l.dirty) {
      l.dirty = true
    }
    return l
  })
  let loaded = count
  return () => {
    if (finishedLoading) {
      return
    }
    if (loaded <= 0) {
      throw `loading()'s response function can only be called a maximum of ${count} times`
    }
    loaded--
    loader.update(l => {
      const val = l.m.get(key)
      val <= 1 ? l.m.delete(key) : l.m.set(key, val - 1)
      return l
    })
  }
}

export const onFinishedLoading = callback => {
  if (!process.browser) {
    return
  }
  let unsubscribe = () => {
    console.error("loading-monitor unsubscribe() called before it is ready!")
  }
  unsubscribe = loader.subscribe(l => {
    //console.log(l.m)
		if (l.m.size === 0 && l.dirty) {
      finishedLoading = true
      callback()
      unsubscribe()
		}
  })
}