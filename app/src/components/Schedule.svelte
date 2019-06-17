<script>
    import { slide } from 'svelte/transition';

    export let schedule = {};
    export let addScheduleHandler = undefined;

    if (!schedule.data.name) {
        schedule.data.name = '';
    }
    
    function open(event) {
        schedule.open = true;
    }

    function close(event) {
        schedule.open = false;
    }

    function save(event) {
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
    p {
        margin: 0 0 0.5rem 0;
    }
    button {
        cursor: pointer;
    }
</style>

<section>
    <header>
        <h2 data-test=schedule-name>
            {#if schedule.editID}
                new schedule [unsaved]
            {:else}
                {`${schedule.data.frequency} schedule`}
            {/if}
        </h2>
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
            <p>[schedule details]</p>
            {#if schedule.editID}
                editing
            {/if}
        </div>
    {/if}
</section>