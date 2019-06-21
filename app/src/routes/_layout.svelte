<script>
	import Nav from '../components/Nav.svelte';
	import { onMount, tick } from 'svelte';

	export let segment;

	let testID = 'loading';
	onMount(() => {
		tick().then(() => {
			testID = 'loaded';
		});
	});

</script>

<style>
	main {
		position: relative;
		max-width: 56em;
		background-color: white;
		padding: 2em;
		margin: 0 auto;
		box-sizing: border-box;
	}
</style>

<Nav {segment}/>

<main data-test={testID}>
	<slot></slot>
</main>