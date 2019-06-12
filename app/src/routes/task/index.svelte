<script context="module">
	import Task from "../../components/Task.svelte"
	
	export async function preload({ params, query }) {
		return this.fetch(`task.json`).then(async r => {
			if (r.status === 200) {
				return r.json()
			} else {
				throw {
					message: "Sorry, there was a problem retrieving tasks",
					r: await r.json()
				};
			}
		}).then(tasks => {
				console.log(tasks);
			return { tasks };
		}).catch(taskError => {
			return { taskError }
		})
	}
</script>

<script>
	export let taskError = undefined;
	export let tasks = {};
	
	if (taskError !== undefined) {
		console.error(taskError);
	}

	const taskArray = Object.values(tasks);
</script>

<style>
	ul {
		list-style: none;
	}
	li {
		padding-bottom: 1px;
	}
	.error {
		color: rgb(199, 25, 60);
	}
	.emptyMessage {
		color: #4d4d4d;
	}
	header h1 {
        display: inline;
	}
	header button {
		float: right;
	}
</style>

<svelte:head>
	<title>Scheduled Tasks - Tasks</title>
</svelte:head>

<section class="tasks">
	<header>
		<h1>Tasks</h1>
		<button>new task</button>
	</header>
	<div class="content">
		{#if taskError !== undefined}
			<p class="error">{taskError.message}</p>
		{/if}

		{#if taskArray.length === 0}
			<p class="emptyMessage">No tasks found</p>
		{:else}
			<ul>
				{#each taskArray as task}
					<li><Task {task}/></li>
				{/each}
			</ul>
		{/if}
	</div>
</section>

<h1>Completed</h1>
<ul>
	<li>Completed Task 3</li>
	<li>Completed Task 4</li>
</ul>
