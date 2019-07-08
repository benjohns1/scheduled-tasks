import { writable } from 'svelte/store'

const loader = writable({
  m: new Map(),
  i: 0,
})

let finishedLoading = false

export const loading = (count = 1) => {
  if (finishedLoading || !process.browser) {
    return () => {}
  }
  let key = undefined
  loader.update(l => {
    key = l.i++
    l.m.set(key, count)
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
  let unsubscribe
	unsubscribe = loader.subscribe(l => {
		if (l.m.size === 0) {
      finishedLoading = true
      callback()
      unsubscribe()
		}
	})
}