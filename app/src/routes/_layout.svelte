<script>
	import Login from '../components/Login.svelte'
	import Nav from '../components/Nav.svelte'
	import { onMount, tick } from 'svelte'
	import { loading, onFinishedLoading } from './../loading-monitor'

	export let segment

	let testID = 'loading'
	onFinishedLoading(() => testID = 'loaded')
	const loaded = loading('layout')
	onMount(() => tick().then(() => loaded()))

</script>

<style>
	.user-status {
		float: right;
		padding: 0 0 1rem 1rem;
	}
	.title {
		float: left;
		padding: 0.25rem 1rem 0 0;
	}

	.container {
		padding: 1rem 0;
	}

	main {
		padding: 2rem 0;
	}
</style>

<div class=container data-test={testID}>
	<div class=title><a href="/"><img src="/icon/favicon-32x32.png" width=32 height=32 alt="Scheduled Tasks" /></div>
	<div class=user-status><Login/></div>
	<Nav {segment}/>

	<main>
		<slot></slot>
	</main>
</div>
