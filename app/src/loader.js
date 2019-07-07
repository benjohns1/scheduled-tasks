import { writable } from 'svelte/store'

export const loader = writable({
  loading: 0,
  loaded: 0,
})

export const addLoader = () => {
  loader.update(l => {
    l.loading += 1
    return l
  })
  let loaded = false
  return () => {
    if (loaded) {
      return
    }
    loaded = true
    loader.update(l => {
      l.loaded += 1
      return l
    })
  }
}

export const onFinishedLoading = callback => {
  let unsubscribe
	unsubscribe = loader.subscribe(({loading, loaded}) => {
		if (loading !== 0 && loading === loaded) {
      callback()
      unsubscribe()
		}
	})
}