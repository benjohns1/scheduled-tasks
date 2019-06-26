<script>
    import { slide } from 'svelte/transition';
	import Task from "./Task.svelte";
	import Button from "./Button.svelte";

    export let schedule = {};
    export let addScheduleHandler = undefined;
    export let deleteScheduleHandler = undefined;

    let tasks = [];

    let ui = {};

    $: {
        if (schedule.data.id !== ui.id || schedule.editID !== ui.editID) {
            ui = {
                id: schedule.data.id,
                editID: schedule.editID,
                minuteMax: 59,
                atMinutes: formatAtMinutes(),
                hourMax: 23,
                atHours: formatAtHours(),
                currTaskEditID: 1
            };
            setAddTaskHandler();

            if (schedule.data.tasks) {
                tasks = schedule.data.tasks.map(t => {
                    return {
                        data: t,
                        open: false
                    };
                });
            }
        }

        ui.name = (() => {
            let frequency = '[unknown]'
            
            const interval = schedule.data.interval !== 1 ? `${schedule.data.interval} ` : '';
            let times = '';
            switch (schedule.data.frequency) {
                case 'Hour':
                    frequency = schedule.data.interval === 1 ? 'hour' : 'hours';
                    times = (schedule.data.atMinutes.length === 1 && schedule.data.atMinutes[0] === 0) ? '' : ` at ${schedule.data.atMinutes.map(m => `${m > 9 ? m : '0' + m}`).join(', ')} minutes`;
                    break;
                case 'Day':
                    frequency = schedule.data.interval === 1 ? 'day' : 'days';
                    times = schedule.data.atHours ? ` at ${schedule.data.atHours.map(h => {
                        return schedule.data.atMinutes.map(m => `${h > 9 ? h : '0' + h}:${m > 9 ? m : '0' + m}`).join(', ');
                    }).join(', ')}` : '';
                    break;
                case 'Week':
                    frequency = schedule.data.interval === 1 ? 'week' : 'weeks';
                    break;
                case 'Month':
                    frequency = schedule.data.interval === 1 ? 'month' : 'months';
                    break;
            }
            return `every ${interval}${frequency}${times}`;
        })();

        setIntervalMax();
    }

    function setIntervalMax() {
        ui.intervalMax = (() => {
            switch (schedule.data.frequency) {
                case 'Hour':
                    return 24;
                case 'Day':
                    return 365;
                case 'Week':
                    return 52;
                case 'Month':
                    return 12;
            }
            return undefined;
        })();
    }

    function formatAtMinutes() {
        return schedule.data.atMinutes ? schedule.data.atMinutes.join(',') : '';
    }

    function formatAtHours() {
        return schedule.data.atHours ? schedule.data.atHours.join(',') : '';
    }

    function frequencyUpdated() {
        validateAll();
    }

    function validateAll() {
        setIntervalMax();

        validateInterval();
        validateOffset();
        validateMinutes();
        validateTasks();
    }

    function validateTasks() {
        // Remove task duplicates
        tasks = tasks.filter((t, i) => {
            for (let j = i + 1; j < tasks.length; j++) {
                if (t.data.name === tasks[j].data.name
                    && t.data.description === tasks[j].data.description) {
                    return false;
                }
            }
            return true;
        });
    }
    
    function validateMinutes() {
        let atMinutes = (ui.atMinutes || '').split(',').reduce((arr, val) => {
            const intVal = parseInt(val);
            const clampedVal = (() => {
                if (intVal < 0) {
                    return 0;
                }
                if (intVal > ui.minuteMax) {
                    return intVal % (ui.minuteMax + 1);
                }
                if (intVal >= 0 && intVal <= ui.minuteMax) {
                    return intVal;
                }
                return undefined
            })();

            if (clampedVal === undefined) {
                return arr;
            }
            if (arr.indexOf(clampedVal) !== -1) {
                return arr;
            }
            return [...arr, clampedVal];
        }, []);
        atMinutes.sort((a, b) => a - b);
        schedule.data.atMinutes = atMinutes;
        ui.atMinutes = formatAtMinutes();
    }

    function validateHours() {
        let atHours = (ui.atHours || '').split(',').reduce((arr, val) => {
            const intVal = parseInt(val);
            const clampedVal = (() => {
                if (intVal < 0) {
                    return 0;
                }
                if (intVal > ui.hourMax) {
                    return intVal % (ui.hourMax + 1);
                }
                if (intVal >= 0 && intVal <= ui.hourMax) {
                    return intVal;
                }
                return undefined
            })();

            if (clampedVal === undefined) {
                return arr;
            }
            if (arr.indexOf(clampedVal) !== -1) {
                return arr;
            }
            return [...arr, clampedVal];
        }, []);
        atHours.sort((a, b) => a - b);
        schedule.data.atHours = atHours;
        ui.atHours = formatAtHours();
    }

    function validateInterval() {
        schedule.data.interval = Math.max(1, Math.min(ui.intervalMax, parseInt(schedule.data.interval)));
    }

    function validateOffset() {
        schedule.data.offset = Math.max(0, Math.min(ui.intervalMax, parseInt(schedule.data.offset)));
    }
    
    function open() {
        schedule.open = true;
    }

    function close() {
        schedule.open = false;
    }

    function save() {
        validateAll();
        if (addScheduleHandler) {
            schedule.data.tasks = tasks.map(t => t.data);
            addScheduleHandler(schedule);
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
		}, ...tasks];
    }

    let addTaskHandler = undefined;
    
	function addTask(taskEditID, taskData) {
		return fetch(`schedule/${schedule.data.id}/task.json`, { method: "POST", headers: {'Content-Type': 'application/json'}, body: JSON.stringify(taskData)}).then(r => {
            if (r.status === 201) {
				tasks = [{
					data: taskData,
					open: true
				}, ...(tasks.filter(t => t.editID !== taskEditID))];
            } else {
                console.error(r);
            }
		}).catch(err => {
			console.error(err);
		});
    }

    function setAddTaskHandler() {
        addTaskHandler = schedule.editID ? undefined : addTask;
    }
    setAddTaskHandler();

    function togglePause() {
        const pause = schedule.data.paused ? 'pause' : 'unpause';
		return fetch(`schedule/${schedule.data.id}/${pause}.json`, { method: "PUT", headers: {'Content-Type': 'application/json'}}).then(r => {
            if (r.status !== 204) {
                console.error(r);
            }
		}).catch(err => {
			console.error(err);
        });
    }

    function deleteSchedule() {
        if (deleteScheduleHandler) {
            deleteScheduleHandler(schedule);
        }
    }

</script>

<style>
    header h3 {
        display: inline;
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
                {#if addScheduleHandler && schedule.editID}
                    <Button on:click={save} test=save-button style=success>save</Button>
                {/if}
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
                            <input id='schedulePaused' type=checkbox class=custom-control-input data-test=paused-toggle bind:checked={schedule.data.paused}>
                        {:else}
                            <input id='schedulePaused' type=checkbox class=custom-control-input data-test=paused-toggle bind:checked={schedule.data.paused} on:change={togglePause}>
                        {/if}
                        <label for='schedulePaused' class='custom-control-label'>Pause ({schedule.data.paused ? 'on' : 'off'})</label>
                    </div>
                </div>
                <div class='form-group row'>
                    {#if schedule.editID}
                        <label for='scheduleFrequency' class='col-sm-2 col-form-label'>Frequency:</label>
                        <div class='col-sm-10'><select id='scheduleFrequency' class=form-control data-test=schedule-frequency-input bind:value={schedule.data.frequency} on:change={frequencyUpdated}>
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
                        <label for='scheduleInterval' class='col-sm-2 col-form-label'>Interval:</label>
                        <div class='col-sm-10'>
                            <input id='scheduleInterval' class=form-control type=number data-test=schedule-interval-input bind:value={schedule.data.interval} min=1 max={ui.intervalMax} on:blur={validateInterval} on:focus={validateInterval}>
                            <small class='form-text text-muted'>(1 - {ui.intervalMax})</small>
                        </div>
                    {:else}
                        <span class='col-sm-2'>Interval:</span> <span class='col-sm-10' data-test=schedule-interval>{schedule.data.interval}</span>
                    {/if}
                </div>
                <div class='form-group row'>
                    {#if schedule.editID}
                        <label for='scheduleOffset' class='col-sm-2 col-form-label'>Offset:</label>
                        <div class='col-sm-10'>
                            <input id='scheduleOffset' class=form-control type=number data-test=schedule-offset-input bind:value={schedule.data.offset} min=0 max={ui.intervalMax} on:blur={validateOffset} on:focus={validateOffset}>
                            <small class='form-text text-muted'>(0 - {ui.intervalMax})</small>
                        </div>
                    {:else}
                        <span class='col-sm-2'>Offset:</span> <span class='col-sm-10' data-test=schedule-offset>{schedule.data.offset}</span>
                    {/if}
                </div>
                {#if schedule.data.frequency === 'Day'}
                    <div class='form-group row'>
                        {#if schedule.editID}
                            <label for='scheduleAtHours' class='col-sm-2 col-form-label'>At hours:</label>
                            <div class='col-sm-10'>
                                <input id='scheduleAtHours' class=form-control type=text data-test=schedule-at-hours-input bind:value={ui.atHours} on:blur={validateHours} on:focus={validateHours}>
                                <small class='form-text text-muted'>(comma-separated, 0 - {ui.minuteMax})</small>
                            </div>
                        {:else}
                            <span class='col-sm-2'>At hours:</span> <span class='col-sm-10' data-test=schedule-at-hours>{ui.atHours}</span>
                        {/if}
                    </div>
                {/if}
                <div class='form-group row'>
                    {#if schedule.editID}
                        <label for='scheduleAtMinutes' class='col-sm-2 col-form-label'>At minutes:</label>
                        <div class='col-sm-10'>
                            <input id='scheduleAtMinutes' class=form-control type=text data-test=schedule-at-minutes-input bind:value={ui.atMinutes} on:blur={validateMinutes} on:focus={validateMinutes}>
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
                <Button on:click={deleteSchedule} test=delete-schedule-button style=outline-danger>delete schedule</Button>
            </div>
        {/if}
    </div>
</section>