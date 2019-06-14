<script context="module">
	import Task from "../../components/Task.svelte"
	
	export function preload({ params, query }) {
		return this.fetch(`task.json`).then(async r => {
			if (r.status === 200) {
				return r.json();
			} else {
				throw {
					message: "Sorry, there was a problem retrieving tasks",
					r: await r.json()
				};
			}
		}).then(tasks => {
			return { tasks: Object.values(tasks).reverse().map(t => {
				return {
					data: t,
					open: false
				}
			}) };
		}).catch(taskError => {
			return { taskError }
		});
	}
</script>

<script>
	export let taskError = undefined;
	export let tasks = [];

	let editID = 1;

	const addTask = (taskEditID, taskData) => {
		return fetch(`task.json`, { method: "POST", headers: {'Content-Type': 'application/json'}, body: JSON.stringify(taskData)}).then(r => {
			r.json().then(({ id }) => {
				taskData.id = id;
				tasks = [{
					data: taskData,
					open: true
				}, ...(tasks.filter(t => t.edit !== taskEditID))];
			});
		}).catch(err => {
			console.error(err);
		});
	}
	
	if (taskError !== undefined) {
		console.error(taskError);
	}

	function newTask() {
		tasks = [{
			data: {
				name: "new task",
				description: ""
			},
			editID: editID++,
			open: true
		}, ...tasks];
		console.log(editID);
	}
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

<section class='tasks'>
	<header>
		<h1>Tasks</h1>
		<button on:click={newTask} data-test='new-task-button'>new task</button>
	</header>
	<div class='content'>
		{#if taskError !== undefined}
			<p class="error">{taskError.message}</p>
		{/if}

		{#if tasks.length === 0}
			<p class='emptyMessage'>No tasks found</p>
		{:else}
			<ul>
				{#each tasks as task}
					<li data-test='task-item'><Task task={task.data} editing={task.editID} opened={task.open} addTaskHandler={addTask}/></li>
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
