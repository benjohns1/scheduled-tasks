<script>
    import { slide } from 'svelte/transition';

    export let schedule = {};
    export let addScheduleHandler = undefined;
    let ui = {
        intervalMax: 0,
        atMinutes: '0',
        name: '[new unsaved schedule]'
    };

    $: ui.name = (() => {
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
        let interval = schedule.data.interval !== 1 ? `${schedule.data.interval} ` : '';
        return `every ${interval}${frequency} at ${schedule.data.atMinutes.map(m => `${m > 9 ? m : '0' + m}`).join(', ')} minutes`;
    })();

    $: ui.intervalMax = (() => {
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

    function showMinutes() {
        ui.atMinutes = schedule.data.atMinutes.join(',');
    }
    
    function validateMinutes() {
        schedule.data.atMinutes = (ui.atMinutes || '').split(',').reduce((arr, val) => {
            const intVal = parseInt(val);
            if (intVal < 0) {
                return [...arr, 0];
            }
            if (intVal > 59) {
                return [...arr, val % 60];
            }
            if (intVal >= 0 && intVal < 60) {
                return [...arr, intVal];
            }
            return arr;
        }, []).sort();
        showMinutes();
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
    showMinutes();
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
                        <option value='Hour'>Hourly</option>
                        <option value='Day'>Daily</option>
                        <option value='Week'>Weekly</option>
                        <option value='Month'>Monthly</option>
                    </select>
                {:else}
                    <span data-test=schedule-frequency id='schedule-frequency'>{schedule.data.frequency}</span>
                {/if}
                </label>
            </div>
            <div>
                <label><span>Interval:</span>
                {#if schedule.editID}
                    <input type=number data-test=schedule-interval-input bind:value={schedule.data.interval} min=1 max={ui.intervalMax}>
                {:else}
                    <span data-test=schedule-interval>{schedule.data.interval}</span>
                {/if}
                </label>
            </div>
            <div>
                <label><span>Offset:</span>
                {#if schedule.editID}
                    <input type=number data-test=schedule-offset-input bind:value={schedule.data.offset} min=0 max={ui.intervalMax}>
                {:else}
                    <span data-test=schedule-offset>{schedule.data.offset}</span>
                {/if}
                </label>
            </div>
            <div>
                <label><span>At minutes:</span>
                {#if schedule.editID}
                    <input type=text data-test=schedule-at-minutes-input bind:value={ui.atMinutes} on:blur={validateMinutes} on:focus={validateMinutes}> (comma-separated)
                {:else}
                    <span data-test=schedule-at-minutes>{ui.atMinutes}</span>
                {/if}
                </label>
            </div>
        </div>
    {/if}
</section>