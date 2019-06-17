<script context="module">
	import Schedule from "../../components/Schedule.svelte"
	
	export function preload({ params, query }) {
		return this.fetch(`schedule.json`).then(r => {
			return r.json().then(data => {
				if (r.status !== 200) {
					console.error(data, r);
					return { scheduleError: 'Sorry, there was a problem retrieving schedules' + data.error ? `: ${data.error}` : ''}
				}

				const schedules = Object.values(data).reverse().map(s => {
					return {
						data: s,
						open: false
					}
				});
				return { schedules }
			});
		}).catch(err => {
			console.error(err);
			return { scheduleError: `Sorry, there was a problem retrieving schedules: ${err.message}` }
		});
	}
</script>

<script>
	export let scheduleError = undefined;
	export let schedules = [];

	let editID = 1;

	const addSchedule = schedule => {
		const editID = schedule.editID;
		return fetch(`schedule.json`, { method: "POST", headers: {'Content-Type': 'application/json'}, body: JSON.stringify(schedule.data)}).then(r => {
			r.json().then(data => {
				if (r.status !== 201) {
					console.error(data, r);
					return;
				}
				schedule.data.id = data.id;
				schedules = [schedule, ...(schedules.filter(s => s.editID !== editID))];
				schedule.editID = undefined;
			});
		}).catch(err => {
			console.error(err);
		});
	}

	function newSchedule() {
		schedules = [{
			data: {
				frequency: "hourly",
				atMinutes: [0,15,30,45]
			},
			editID: editID++,
			open: true
		}, ...schedules];	
	}
</script>

<style>
	ul {
		list-style: none;
	}
	li {
		padding-bottom: 1px;
		clear: both;
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
	<title>Scheduled Tasks - Schedules</title>
</svelte:head>

<section data-test=schedules>
	<header>
		<h1>Schedules</h1>
		<button on:click={newSchedule} data-test=new-schedule-button>new schedule</button>
	</header>
	{#if scheduleError !== undefined}
		<p class='error'>{scheduleError}</p>
	{/if}

	{#if schedules.length === 0}
		<p class='emptyMessage'>No schedules found</p>
	{:else}
		<ul>
			{#each schedules as schedule}
				<li data-test=schedule-item><Schedule schedule={schedule} addScheduleHandler={addSchedule}/></li>
			{/each}
		</ul>
	{/if}
</section>