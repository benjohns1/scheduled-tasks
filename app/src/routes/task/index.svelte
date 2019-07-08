<script context="module">
	import Task from "../../components/Task.svelte"
	import Button from "../../components/Button.svelte"
	import { withJsonAndAuth } from "../../api/default.headers"
	import { loading } from './../../loading-monitor'

	const loaded = loading()
	
	export async function preload(page, session ) {
		const result = await this.fetch(`task.json`, withJsonAndAuth(session ? session.token : null)).then(async r => {
			if (r.status === 200) {
				return r.json()
			} else {
				throw {
					message: "Sorry, there was a problem retrieving tasks",
					r: await r.json()
				}
			}
		}).then(tasks => {
			const allTasks = Object.values(tasks).reverse().map(t => {
				return {
					data: t,
					open: false
				}
			})
			return { tasks: allTasks.filter(t => !t.data.completedTime), completedTasks: allTasks.filter(t => t.data.completedTime) }
		}).catch(taskError => {
			return { taskError }
		})
	
		loaded()
		return result
	}
</script>

<script>
    import { stores } from '@sapper/app'
	const { session } = stores()
	
	export let taskError = undefined
	export let tasks = []
	export let completedTasks = []
	export let completedSuccessMessage = undefined

	let editID = 1

	const addTask = (taskEditID, taskData) => {
		return fetch(`task.json`, withJsonAndAuth($session.token, { method: "POST", body: JSON.stringify(taskData)})).then(r => {
			r.json().then(({ id }) => {
				taskData.id = id
				tasks = [{
					data: taskData,
					open: true
				}, ...(tasks.filter(t => t.editID !== taskEditID))]
			})
		}).catch(err => {
			console.error(err)
		})
	}

	const completeTask = taskID => {
		return fetch(`task/${taskID}/complete.json`, withJsonAndAuth($session.token, { method: "PUT" })).then(r => {
			completedTasks = [...(tasks.filter(t => t.data.id === taskID).map(t => {
				t.data.completedTime = 'just now'
				t.open = false
				return t
			})), ...completedTasks]
			tasks = [...(tasks.filter(t => t.data.id !== taskID))]
		}).catch(err => {
			console.error(err)
		})
	}
	
	if (taskError !== undefined) {
		console.error(taskError)
	}

	function newTask() {
		tasks = [{
			data: {
				name: "new task",
				description: ""
			},
			editID: editID++,
			open: true
		}, ...tasks]
	}

	function clearTasks() {
		const prevCompletedTasks = completedTasks.slice(0)
		completedTasks = []
		return fetch(`task/clear.json`, withJsonAndAuth($session.token, { method: "POST" } )).then(r => {
			r.json().then(({ count, message }) => {
				completedSuccessMessage = `${message}${count ? ` (${count})` : ''}`
			})
		}).catch(err => {
			completedTasks = prevCompletedTasks
			console.error(err)
		})
	}
</script>

<style>
	ul {
		list-style: none;
        margin: 1px 0;
        padding: 0;
	}
	li {
		padding-bottom: 1px;
		clear: both;
	}
	.error {
		color: rgb(199, 25, 60);
	}
	.successMessage {
		color: #2dc066;
	}
	.emptyMessage {
		color: #4d4d4d;
	}
	header {
		margin-bottom: 0.5rem;
	}
	header h1 {
        display: inline;
	}
    header:after {
        content: "";
        clear: both;
        display: table;
	}
	section {
		margin-bottom: 1rem;
	}
</style>

<svelte:head>
	<title>Scheduled Tasks - Tasks</title>
</svelte:head>

<section data-test=tasks>
	<header>
		<h1>Tasks</h1>
		<Button on:click={newTask} test=new-task-button classes=right style=success>new task</Button>
	</header>
	{#if taskError !== undefined}
		<p class="error">{taskError.message}</p>
	{/if}

	{#if tasks.length === 0}
		<p class='emptyMessage'>No tasks found</p>
	{:else}
		<ul>
			{#each tasks as task}
				<li data-test=task-item><Task task={task.data} editing={task.editID} opened={task.open} addTaskHandler={addTask} completeTaskHandler={completeTask}/></li>
			{/each}
		</ul>
	{/if}
</section>

<section data-test=completed-tasks>
	<header>
		<h1>Completed</h1>
		<Button on:click={clearTasks} test=clear-tasks-button classes=right style=outline-danger>clear all completed tasks</Button>
	</header>
	{#if completedSuccessMessage !== undefined}
		<p class="successMessage" data-test=completed-success-message>{completedSuccessMessage}</p>
	{/if}

	{#if completedTasks.length === 0}
		<p class='emptyMessage' data-test=completed-empty-message>No completed tasks</p>
	{:else}
		<ul>
			{#each completedTasks as task}
				<li data-test=completed-task-item><Task task={task.data} opened={task.open}/></li>
			{/each}
		</ul>
	{/if}
</section>