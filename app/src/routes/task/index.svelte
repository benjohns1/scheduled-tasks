<script context="module">
	export async function preload({ params, query }) {
		return this.fetch(`task.json`).then(r => r.json()).then(tasks => {
			return { tasks };
		});
	}
</script>

<script>
	export let tasks;
</script>

<svelte:head>
	<title>Scheduled Tasks - Tasks</title>
</svelte:head>

<h1>Tasks</h1>
<ul>
	{#each Object.values(tasks) as task}
		<li><a href="task/{task.id}">{task.name}</a></li>
	{/each}
</ul>

<h1>Completed</h1>
<ul>
	<li>Completed Task 3</li>
	<li>Completed Task 4</li>
</ul>
