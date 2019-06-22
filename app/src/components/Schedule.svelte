<script>
    import { slide } from 'svelte/transition';
	import Task from "./Task.svelte";

    export let schedule = {};
    export let addScheduleHandler = undefined;

    let tasks = [];

    let ui = {};

    $: {
        if (schedule.data.id !== ui.id || schedule.editID !== ui.editID) {
            ui = {
                id: schedule.data.id,
                editID: schedule.editID,
                minuteMax: 59,
                atMinutes: formatAtMinutes(),
                currTaskEditID: 1
            };
            setAddTaskHandler();

            if (schedule.data && schedule.data.tasks) {
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
            switch (schedule.data.frequency) {
                case 'Hour':
                    frequency = schedule.data.interval === 1 ? 'hour' : 'hours';
                    break;
                case 'Day':
                    frequency = schedule.data.interval === 1 ? 'day' : 'days';
                    break;
                case 'Week':
                    frequency = schedule.data.interval === 1 ? 'week' : 'weeks';
                    break;
                case 'Month':
                    frequency = schedule.data.interval === 1 ? 'month' : 'months';
                    break;
            }
            const interval = schedule.data.interval !== 1 ? `${schedule.data.interval} ` : '';
            const minutes =  (schedule.data.atMinutes.length === 1 && schedule.data.atMinutes[0] === 0) ? '' : ` at ${schedule.data.atMinutes.map(m => `${m > 9 ? m : '0' + m}`).join(', ')} minutes`;
            return `every ${interval}${frequency}${minutes}`;
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
        return schedule.data.atMinutes.join(',');
    }

    function validateAll() {
        setIntervalMax();

        validateInterval();
        validateOffset();
        validateMinutes();
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
				name: "recurring task",
				description: ""
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

</script>

<style>
    header {
        background-color: #ddd;
        padding: 0.4rem 1rem;
    }
    header h2 {
        display: inline;
    }
    span.right {
        float: right;
        margin-left: 1rem;
    }
    section {
        background-color: #eee;
    }
	.emptyMessage {
		color: #4d4d4d;
	}
    .panel {
        margin: 0;
        padding: 0.4rem 1rem;
    }
    button {
        cursor: pointer;
    }
    label span:nth-child(1) {
        display: inline-block;
        width: 6.5rem;
    }
	ul {
		list-style: none;
        margin: 1px 0;
        padding: 0;
	}
	li {
		padding-bottom: 1px;
		clear: both;
	}
</style>

<section>
    <header>
        <h2 data-test=schedule-name>{ui.name}</h2>
        <span class=right>
            {#if addScheduleHandler && schedule.editID}
                <button on:click={save} data-test=save-button>save</button>
            {/if}
            {#if schedule.open}
                <button on:click={close} data-test=close-button>v</button>
            {:else}
                <button on:click={open} data-test=open-button>&gt;</button>
            {/if}
        </span>
    </header>
    {#if schedule.open}
        <div class='panel' transition:slide='{{ duration: 100 }}'>
            <div>
                <label for='schedule-frequency'><span>Frequency:</span>
                {#if schedule.editID}
                    <select id='schedule-frequency' data-test=schedule-frequency-input bind:value={schedule.data.frequency} on:change={validateAll}>
                        <option value='Hour'>Hour</option>
                        <option value='Day'>Day</option>
                        <option value='Week'>Week</option>
                        <option value='Month'>Month</option>
                    </select>
                {:else}
                    <span data-test=schedule-frequency id='schedule-frequency'>{schedule.data.frequency}</span>
                {/if}
                </label>
            </div>
            <div>
                <label><span>Interval:</span>
                {#if schedule.editID}
                    <input type=number data-test=schedule-interval-input bind:value={schedule.data.interval} min=1 max={ui.intervalMax} on:blur={validateInterval} on:focus={validateInterval}> (1 - {ui.intervalMax})
                {:else}
                    <span data-test=schedule-interval>{schedule.data.interval}</span>
                {/if}
                </label>
            </div>
            <div>
                <label><span>Offset:</span>
                {#if schedule.editID}
                    <input type=number data-test=schedule-offset-input bind:value={schedule.data.offset} min=0 max={ui.intervalMax} on:blur={validateOffset} on:focus={validateOffset}> (0 - {ui.intervalMax})
                {:else}
                    <span data-test=schedule-offset>{schedule.data.offset}</span>
                {/if}
                </label>
            </div>
            <div>
                <label><span>At minutes:</span>
                {#if schedule.editID}
                    <input type=text data-test=schedule-at-minutes-input bind:value={ui.atMinutes} on:blur={validateMinutes} on:focus={validateMinutes}> (comma-separated, 0 - {ui.minuteMax})
                {:else}
                    <span data-test=schedule-at-minutes>{ui.atMinutes}</span>
                {/if}
                </label>
            </div>
            <div>
                <label><span>Paused:</span>
                {#if schedule.editID}
                    <input type=checkbox data-test=paused-toggle bind:checked={schedule.data.paused}>
                {:else}
                    <input type=checkbox data-test=paused-toggle bind:checked={schedule.data.paused} on:change={togglePause}>
                {/if}
                </label>
            </div>
            <div>
                <span class='right'><button on:click={newTask} data-test=new-task>new task</button></span>
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
        </div>
    {/if}
</section>