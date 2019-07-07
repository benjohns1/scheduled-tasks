<script>
	import Login from '../components/Login.svelte'
	import Nav from '../components/Nav.svelte'
	import { onMount, tick } from 'svelte'

	export let segment

	let testID = 'loading'
	onMount(() => {
		tick().then(() => {
			testID = 'loaded'
		})
	})

</script>

<style>
	.user-status {
		float: right;
		padding: 0 0 2rem 1rem;
	}

	.container {
		padding: 1rem 0;
	}

	main {
		padding: 2rem 0;
	}
</style>

<div class=container data-test={testID}>
	<div class=user-status><Login/></div>
	<Nav {segment}/>

	<main>
		<slot></slot>
	</main>
</div>