<script>
    import { slide } from 'svelte/transition';

    export let schedule = {};
    export let addScheduleHandler = undefined;

    let ui = {};

    $: {
        if (schedule.data.id !== ui.id || schedule.editID !== ui.editID) {
            ui = {
                id: schedule.data.id,
                editID: schedule.editID,
                minuteMax: 59,
                atMinutes: formatAtMinutes()
            };
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
    
    function open() {
        schedule.open = true;
    }

    function close() {
        schedule.open = false;
    }

    function save() {
        validateMinutes();
        if (addScheduleHandler) {
            addScheduleHandler(schedule);
        }
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
                    <select id='schedule-frequency' data-test=schedule-frequency-input bind:value={schedule.data.frequency}>
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
                    <input type=number data-test=schedule-interval-input bind:value={schedule.data.interval} min=1 max={ui.intervalMax}> (0 - {ui.intervalMax})
                {:else}
                    <span data-test=schedule-interval>{schedule.data.interval}</span>
                {/if}
                </label>
            </div>
            <div>
                <label><span>Offset:</span>
                {#if schedule.editID}
                    <input type=number data-test=schedule-offset-input bind:value={schedule.data.offset} min=0 max={ui.intervalMax}> (0 - {ui.intervalMax})
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
        </div>
    {/if}
</section>