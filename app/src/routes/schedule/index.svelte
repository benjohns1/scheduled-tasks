<script context="module">
	import Schedule from "../../components/Schedule.svelte"
	import Button from "../../components/Button.svelte"
	import { withJsonAndAuth } from "../../api/default.headers"
	import { loading } from './../../loading-monitor'
	import Login from '../../components/Login.svelte'

	export async function preload(page, session) {
		const schedulesFetched = loading('fetchSchedules')
		const result = await this.fetch(`schedule.json`, withJsonAndAuth(session)).then(async r => {
			if (r.status !== 200) {
				throw {
					message: "Sorry, there was a problem retrieving schedules",
					r: await r.json()
				}
			} else {
				const schedules = Object.values(await r.json()).reverse().map(s => {
					return {
						data: s,
						open: false
					}
				})
				return { schedules }
			}
		}).catch(scheduleError => {
			console.error(scheduleError)
			return { scheduleError }
		})
		schedulesFetched()
		return result
	}
</script>

<script>
	import { onMount, tick } from 'svelte'
    import { stores } from '@sapper/app'
	const { session } = stores()

	export let scheduleError = undefined
	export let schedules = []

	let editID = 1
	const routeLoaded = loading('scheduleRoute')

	const addSchedule = schedule => {
		const editID = schedule.editID
		const postData = (() => {
			const ignoreProps = (() => {
				switch (schedule.data.frequency) {
					case 'Hour':
						return ['atHours', 'onDaysOfMonth', 'onDaysOfWeek']
					case 'Day':
						return ['onDaysOfMonth', 'onDaysOfWeek']
					case 'Week':
						return ['onDaysOfMonth']
					case 'Month':
						return ['onDaysOfWeek']
				}
				return []
			})()
			if (!schedule.data.tasks || schedule.data.tasks.length <= 0) {
				ignoreProps.push('tasks')
			}
			return Object.keys(schedule.data).reduce((acc, key) => {
				if (schedule.data.hasOwnProperty(key) && !ignoreProps.includes(key)) {
					acc[key] = schedule.data[key]
				}
				return acc
			}, {})
		})()
		return fetch(`schedule.json`, withJsonAndAuth($session, { method: "POST", body: JSON.stringify(postData)})).then(r => {
			r.json().then(data => {
				if (r.status !== 201) {
					console.error(data, r)
					return
				}
				schedule.data.id = data.id
				schedules = [schedule, ...(schedules.filter(s => s.editID !== editID))]
				schedule.editID = undefined
			})
		}).catch(err => {
			console.error(err)
		})
	}

    const deleteSchedule = schedule => {
		if (schedule.data.id !== undefined) {
			const id = schedule.data.id
			return fetch(`schedule/${id}.json`, withJsonAndAuth($session, { method: "DELETE" })).then(r => {
				if (r.status !== 204) {
					console.error(r)
					return
				}
				schedules = schedules.filter(s => s.data.id !== id)
			}).catch(err => {
				console.error(err)
			})
		}

		if (schedule.editID !== undefined) {
			const editID = schedule.editID
			schedules = schedules.filter(s => s.editID !== editID)
			return
		}

		console.error('error deleting schedule, no valid ID', schedule)
    }

	function newSchedule() {
		schedules = [{
			data: {
				frequency: "Hour",
				interval: 1,
				offset: 0,
				atMinutes: [0],
				atHours: [0],
				onDaysOfWeek: ['Monday'],
				onDaysOfMonth: [1],
				tasks: []
			},
			editID: editID++,
			open: true
		}, ...schedules]
	}

	onMount(async () => {
		await tick()
		routeLoaded()
	})
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
</style>

<svelte:head>
	<title>Scheduled Tasks - Schedules</title>
</svelte:head>

{#if !$session.auth.isAuthenticated}
	<Login/> to view your schedules
{:else}
	<section data-test=schedules>
		<header>
			<h1>Schedules</h1>
			<Button on:click={newSchedule} test=new-schedule-button classes=right style=success>new schedule</Button>
		</header>
		{#if scheduleError !== undefined}
			<p class='error'>{scheduleError.message || 'Unknown error'}</p>
		{/if}

		{#if schedules.length === 0}
			<p class='emptyMessage'>No schedules found</p>
		{:else}
			<ul>
				{#each schedules as schedule}
					<li data-test=schedule-item><Schedule schedule={schedule} addScheduleHandler={addSchedule} deleteScheduleHandler={deleteSchedule}/></li>
				{/each}
			</ul>
		{/if}
	</section>
{/if}
