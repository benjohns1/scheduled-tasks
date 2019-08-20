<script>
    import { slide } from 'svelte/transition'
	import Task from "./Task.svelte"
	import Button from "./Button.svelte"
	import { withJsonAndAuth } from "../api/default.headers"
    import { stores } from '@sapper/app'
	const { session } = stores()

    export let schedule = {}
    export let addScheduleHandler = undefined
    export let deleteScheduleHandler = undefined

    let tasks = []

    let ui = {}

    $: {
        if (schedule.data.id !== ui.id || schedule.editID !== ui.editID) {
            ui = {
                id: schedule.data.id,
                editID: schedule.editID,
                key: schedule.data.id !== undefined ? `id-${schedule.data.id}` : `edit-${schedule.editID}`,
                minuteMax: 59,
                atMinutes: formatAtMinutes(),
                hourMax: 23,
                atHours: formatAtHours(),
                weekdays: ['Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday'],
                onDaysOfWeek: formatOnDaysOfWeek(),
                dayMax: 31,
                onDaysOfMonth: formatOnDaysOfMonth(),
                currTaskEditID: 1
            }
            setAddTaskHandler()

            if (schedule.data.tasks) {
                tasks = schedule.data.tasks.map(t => {
                    return {
                        data: t,
                        open: false
                    }
                })
            }
        }

        ui.name = (() => {
            let frequency = '[unknown]'

            const formatHoursAndMinutes = () => {
                return schedule.data.atHours ? ` at ${schedule.data.atHours.map(h => {
                    return schedule.data.atMinutes.map(m => `${h > 9 ? h : '0' + h}:${m > 9 ? m : '0' + m}`).join(', ')
                }).join(', ')}` : ''
            }

            const addOrdinalSuffix = (num) => {
                if (![11,12,13].includes(num % 100)) {
                    switch (num % 10) {
                        case 1:
                            return `${num}st`
                        case 2:
                            return `${num}nd`
                        case 3:
                            return `${num}rd`
                    }
                }
                return `${num}th`
            }
            
            const interval = schedule.data.interval !== 1 ? `${schedule.data.interval} ` : ''
            let times = ''
            let isValid = (() => {
                switch (schedule.data.frequency) {
                    case 'Hour':
                        return  schedule.data.atMinutes && schedule.data.atMinutes.length > 0
                    case 'Day':
                        return (schedule.data.atMinutes && schedule.data.atHours) && (schedule.data.atMinutes.length > 0 && schedule.data.atHours.length > 0)
                    case 'Week':
                        return (schedule.data.atMinutes && schedule.data.atHours && schedule.data.onDaysOfWeek) && (schedule.data.atMinutes.length > 0 && schedule.data.atHours.length > 0 && schedule.data.onDaysOfWeek.length > 0)
                    case 'Month':
                        return (schedule.data.atMinutes && schedule.data.atHours && schedule.data.onDaysOfMonth) && (schedule.data.atMinutes.length > 0 && schedule.data.atHours.length > 0 && schedule.data.onDaysOfMonth.length > 0)
                }
                return false
            })()
            if (!isValid) {
                return 'no recurrences scheduled'
            }
            switch (schedule.data.frequency) {
                case 'Hour':
                    frequency = schedule.data.interval === 1 ? 'hour' : 'hours'
                    times = (schedule.data.atMinutes.length === 1 && schedule.data.atMinutes[0] === 0) ? '' : ` at ${schedule.data.atMinutes.map(m => `${m > 9 ? m : '0' + m}`).join(', ')} minutes`
                    break
                case 'Day':
                    frequency = schedule.data.interval === 1 ? 'day' : 'days'
                    times = formatHoursAndMinutes()
                    break
                case 'Week':
                    frequency = schedule.data.interval === 1 ? 'week' : 'weeks'
                    times = ` on ${schedule.data.onDaysOfWeek.join(', ')}${formatHoursAndMinutes()}`
                    break
                case 'Month':
                    frequency = schedule.data.interval === 1 ? 'month' : 'months'
                    times = ` on the ${schedule.data.onDaysOfMonth.map(d => addOrdinalSuffix(d)).join(', ')}${formatHoursAndMinutes()}`
                    break
            }
            return `every ${interval}${frequency}${times}`
        })()

        setIntervalMax()
    }

    function setIntervalMax() {
        ui.intervalMax = (() => {
            switch (schedule.data.frequency) {
                case 'Hour':
                    return 24
                case 'Day':
                    return 365
                case 'Week':
                    return 52
                case 'Month':
                    return 12
            }
            return undefined
        })()
    }

    function formatAtMinutes() {
        return schedule.data.atMinutes ? schedule.data.atMinutes.join(', ') : ''
    }

    function formatAtHours() {
        return schedule.data.atHours ? schedule.data.atHours.join(', ') : ''
    }

    function formatOnDaysOfWeek() {
        if (!ui.weekdays || !schedule.data.onDaysOfWeek) {
            return
        }
        return ui.weekdays.reduce((acc, day) => {
            acc[day] = schedule.data.onDaysOfWeek.includes(day)
            return acc
        }, {})
    }

    function formatOnDaysOfMonth() {
        return schedule.data.onDaysOfMonth ? schedule.data.onDaysOfMonth.join(', ') : ''
    }

    function frequencyUpdated() {
        validateAll()
    }

    function daysOfWeekUpdated() {
        validateAll()
    }

    function validateAll() {
        setIntervalMax()

        validateInterval()
        validateOffset()
        validateMinutes()
        validateHours()
        validateWeekDays()
        validateMonthDays()
        validateTasks()
    }

    function validateTasks() {
        // Remove task duplicates
        tasks = tasks.filter((t, i) => {
            for (let j = i + 1; j < tasks.length; j++) {
                if (t.data.name === tasks[j].data.name
                    && t.data.description === tasks[j].data.description) {
                    return false
                }
            }
            return true
        })
    }
    
    function validateMinutes() {
        schedule.data.atMinutes = commaListToArray(ui.atMinutes, 0, ui.minuteMax)
        ui.atMinutes = formatAtMinutes()
    }

    function validateHours() {
        schedule.data.atHours = commaListToArray(ui.atHours, 0, ui.hourMax)
        ui.atHours = formatAtHours()
    }

    function validateWeekDays() {
        if (!ui.onDaysOfWeek) {
            return
        }
        schedule.data.onDaysOfWeek = Object.keys(ui.onDaysOfWeek).reduce((acc, day) => {
            if (ui.onDaysOfWeek[day]) {
                acc.push(day)
            }
            return acc
        }, [])
        ui.onDaysOfWeek = formatOnDaysOfWeek()
    }

    function validateMonthDays() {
        schedule.data.onDaysOfMonth = commaListToArray(ui.onDaysOfMonth, 1, ui.dayMax)
        ui.onDaysOfMonth = formatOnDaysOfMonth()
    }

    function commaListToArray(commaList, min, max) {
        let values = (commaList || '').split(',').reduce((arr, val) => {
            const intVal = parseInt(val)
            const clampedVal = Math.min(Math.max(intVal, min), max)

            if (clampedVal === undefined) {
                return arr
            }
            if (arr.indexOf(clampedVal) !== -1) {
                return arr
            }
            return [...arr, clampedVal]
        }, [])
        values.sort((a, b) => a - b)
        return values
    }

    function validateInterval() {
        schedule.data.interval = Math.max(1, Math.min(ui.intervalMax, parseInt(schedule.data.interval)))
    }

    function validateOffset() {
        schedule.data.offset = Math.max(0, Math.min(ui.intervalMax, parseInt(schedule.data.offset)))
    }
    
    function open() {
        schedule.open = true
    }

    function close() {
        schedule.open = false
    }

    function save() {
        validateAll()
        if (addScheduleHandler) {
            schedule.data.tasks = tasks.map(t => t.data)
            addScheduleHandler(schedule)
        }
    }

	function newTask() {
		tasks = [{
			data: {
				name: `recurring task ${tasks.length}`,
				description: ''
			},
			editID: ui.currTaskEditID++,
			open: true
		}, ...tasks]
    }

    let addTaskHandler = undefined
	function addTask(taskEditID, taskData) {
		return fetch(`schedule/${schedule.data.id}/task.json`, withJsonAndAuth($session, { method: "POST", body: JSON.stringify(taskData) })).then(r => {
            if (r.status === 201) {
				tasks = [{
					data: taskData,
					open: true
				}, ...(tasks.filter(t => t.editID !== taskEditID))]
            } else {
                console.error(r)
            }
		}).catch(err => {
			console.error(err)
		})
    }

    function setAddTaskHandler() {
        addTaskHandler = schedule.editID ? undefined : addTask
    }
    setAddTaskHandler()

    function togglePause() {
        const pause = schedule.data.paused ? 'pause' : 'unpause'
		return fetch(`schedule/${schedule.data.id}/${pause}.json`, withJsonAndAuth($session, { method: "PUT" })).then(r => {
            if (r.status !== 204) {
                console.error(r)
            }
		}).catch(err => {
			console.error(err)
        })
    }

    function deleteSchedule() {
        if (deleteScheduleHandler) {
            deleteScheduleHandler(schedule)
        }
    }

</script>

<style>
    header h3 {
        display: inline;
    }
    footer {
        margin-top: 0.5rem;
    }
    .right {
        float: right;
        margin-left: 1rem;
    }
    .clearfix:after,
    header:after {
        content: "";
        clear: both;
        display: table;
    }
	.emptyMessage {
		color: #4d4d4d;
	}
	ul {
		list-style: none;
        margin: 1px 0;
        padding: 0.5rem 0 0 0;
	}
	li {
		padding-bottom: 1px;
		clear: both;
	}
    .card-text {
        padding-top: 0.5rem;
    }
    .custom-switch input[type=checkbox] {
        position: relative;
    }
</style>

<section class=card>
    <div class=card-body>
        <header>
            <h3 data-test=schedule-name class=card-title>{ui.name}</h3>
            <span class=right>
                {#if schedule.open}
                    <Button on:click={close} test=close-button style=secondary>v</Button>
                {:else}
                    <Button on:click={open} test=open-button style=secondary>&gt;</Button>
                {/if}
            </span>
        </header>
        {#if schedule.open}
            <div class=card-text transition:slide='{{ duration: 100 }}'>
                <div class='form-group row'>
                    <div class='custom-control custom-switch'>
                        {#if schedule.editID}
                            <input id='schedulePaused-{ui.key}' type=checkbox class=custom-control-input data-test=paused-toggle bind:checked={schedule.data.paused}>
                        {:else}
                            <input id='schedulePaused-{ui.key}' type=checkbox class=custom-control-input data-test=paused-toggle bind:checked={schedule.data.paused} on:change={togglePause}>
                        {/if}
                        <label for='schedulePaused-{ui.key}' class='custom-control-label'>Pause ({schedule.data.paused ? 'on' : 'off'})</label>
                    </div>
                </div>
                <div class='form-group row'>
                    {#if schedule.editID}
                        <label for='scheduleFrequency-{ui.key}' class='col-sm-2 col-form-label'>Frequency:</label>
                        <div class='col-sm-10'><select id='scheduleFrequency-{ui.key}' class=form-control data-test=schedule-frequency-input bind:value={schedule.data.frequency} on:change={frequencyUpdated}>
                            <option value='Hour'>Hour</option>
                            <option value='Day'>Day</option>
                            <option value='Week'>Week</option>
                            <option value='Month'>Month</option>
                        </select></div>
                    {:else}
                        <span class='col-sm-2'>Frequency:</span> <span class='col-sm-10' data-test=schedule-frequency>{schedule.data.frequency}</span>
                    {/if}
                </div>
                <div class='form-group row'>
                    {#if schedule.editID}
                        <label for='scheduleInterval-{ui.key}' class='col-sm-2 col-form-label'>Interval:</label>
                        <div class='col-sm-10'>
                            <input id='scheduleInterval-{ui.key}' class=form-control type=number data-test=schedule-interval-input bind:value={schedule.data.interval} min=1 max={ui.intervalMax} on:blur={validateInterval} on:focus={validateInterval}>
                            <small class='form-text text-muted'>(1 - {ui.intervalMax})</small>
                        </div>
                    {:else}
                        <span class='col-sm-2'>Interval:</span> <span class='col-sm-10' data-test=schedule-interval>{schedule.data.interval}</span>
                    {/if}
                </div>
                <div class='form-group row'>
                    {#if schedule.editID}
                        <label for='scheduleOffset-{ui.key}' class='col-sm-2 col-form-label'>Offset:</label>
                        <div class='col-sm-10'>
                            <input id='scheduleOffset-{ui.key}' class=form-control type=number data-test=schedule-offset-input bind:value={schedule.data.offset} min=0 max={ui.intervalMax} on:blur={validateOffset} on:focus={validateOffset}>
                            <small class='form-text text-muted'>(0 - {ui.intervalMax})</small>
                        </div>
                    {:else}
                        <span class='col-sm-2'>Offset:</span> <span class='col-sm-10' data-test=schedule-offset>{schedule.data.offset}</span>
                    {/if}
                </div>
                {#if schedule.data.frequency === 'Month'}
                    <div class='form-group row'>
                        {#if schedule.editID}
                            <label for='scheduleOnDaysOfMonth-{ui.key}' class='col-sm-2 col-form-label'>On days:</label>
                            <div class='col-sm-10'>
                                <input id='scheduleOnDaysOfMonth-{ui.key}' class=form-control type=text data-test=schedule-on-days-of-month-input bind:value={ui.onDaysOfMonth} on:blur={validateMonthDays} on:focus={validateMonthDays}>
                                <small class='form-text text-muted'>(comma-separated, 1 - {ui.dayMax})</small>
                            </div>
                        {:else}
                            <span class='col-sm-2'>On days:</span> <span class='col-sm-10' data-test=schedule-on-days-of-month>{ui.onDaysOfMonth}</span>
                        {/if}
                    </div>
                {/if}
                {#if schedule.data.frequency === 'Week'}
                    <div class='form-group row'>
                        <span class='col-sm-2'>On days:</span>
                        {#if schedule.editID}
                            {#each ui.weekdays as day}
                                <div class="form-check form-check-inline">
                                    <input class=form-check-input data-test='schedule-on-days-of-week-input-{day}' type=checkbox id='dayOfWeek{day}-{ui.key}' bind:checked={ui.onDaysOfWeek[day]} on:change={daysOfWeekUpdated}>
                                    <label class=form-check-label for='dayOfWeek{day}-{ui.key}'>{day}</label>
                                </div>
                            {/each}
                        {:else}
                            <span class='col-sm-10' data-test=schedule-on-days-of-week>{schedule.data.onDaysOfWeek ? schedule.data.onDaysOfWeek.join(', ') : '[none]'}</span>
                        {/if}
                    </div>
                {/if}
                {#if schedule.data.frequency !== 'Hour'}
                    <div class='form-group row'>
                        {#if schedule.editID}
                            <label for='scheduleAtHours-{ui.key}' class='col-sm-2 col-form-label'>At hours:</label>
                            <div class='col-sm-10'>
                                <input id='scheduleAtHours-{ui.key}' class=form-control type=text data-test=schedule-at-hours-input bind:value={ui.atHours} on:blur={validateHours} on:focus={validateHours}>
                                <small class='form-text text-muted'>(comma-separated, 0 - {ui.hourMax})</small>
                            </div>
                        {:else}
                            <span class='col-sm-2'>At hours:</span> <span class='col-sm-10' data-test=schedule-at-hours>{ui.atHours}</span>
                        {/if}
                    </div>
                {/if}
                <div class='form-group row'>
                    {#if schedule.editID}
                        <label for='scheduleAtMinutes-{ui.key}' class='col-sm-2 col-form-label'>At minutes:</label>
                        <div class='col-sm-10'>
                            <input id='scheduleAtMinutes-{ui.key}' class=form-control type=text data-test=schedule-at-minutes-input bind:value={ui.atMinutes} on:blur={validateMinutes} on:focus={validateMinutes}>
                            <small class='form-text text-muted'>(comma-separated, 0 - {ui.minuteMax})</small>
                        </div>
                    {:else}
                        <span class='col-sm-2'>At minutes:</span> <span class='col-sm-10' data-test=schedule-at-minutes>{ui.atMinutes}</span>
                    {/if}
                </div>
                <div>
                    <div class=clearfix><Button on:click={newTask} test=new-task classes=right style=success>new task</Button></div>
                    {#if tasks.length === 0}
                        <p class='emptyMessage'>No recurring tasks</p>
                    {:else}
                        <ul>
                            {#each tasks as task}
                                <li data-test=task-item><Task task={task.data} editing={task.editID} opened={task.open} addTaskHandler={addTaskHandler}/></li>
                            {/each}
                        </ul>
                    {/if}
                </div>
                <footer>
                    <Button on:click={deleteSchedule} test=delete-schedule-button style=outline-danger>delete schedule</Button>
                    <div class=right>
                        {#if addScheduleHandler && schedule.editID}
                            <Button on:click={save} test=save-button style=success>save</Button>
                        {/if}
                    </div>
                </footer>
            </div>
        {/if}
    </div>
</section>