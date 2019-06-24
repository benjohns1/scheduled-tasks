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

<div class=container>
	<Nav {segment}/>

	<main data-test={testID}>
		<slot></slot>
	</main>
</div>