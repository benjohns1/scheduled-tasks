<script>
    import { slide } from 'svelte/transition';

    export let task = {};
    export let opened = false;
    export let editing = undefined;
    export let addTaskHandler = undefined;
    export let completeTaskHandler = undefined;

    if (!task.name) {
        task.name = '';
    }
    if (!task.description) {
        task.description = '';
    }
    
    function open() {
        opened = true;
    }

    function close() {
        opened = false;
    }

    function save() {
        if (addTaskHandler) {
            addTaskHandler(editing, task);
        }
    }

    function complete() {
        if (completeTaskHandler && task.id) {
            completeTaskHandler(task.id);
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
    span.left {
        float: left;
        margin-right: 1rem;
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
    textarea {
        width: 100%;
        height: 100%;
    }
</style>

<section>
    <header>
        <h2 data-test=task-name>
            {#if editing}
                <input type=text bind:value={task.name} placeholder='task name' data-test=task-name-input>
            {:else}
                {(task.name || 'task')}
            {/if}
        </h2>
        {#if !editing && !task.completedTime && completeTaskHandler}
            <span class='left'>
                <button on:click={complete} data-test=complete-toggle>done</button>
            </span>
        {/if}
        <span class=right>
            {#if editing && addTaskHandler}
                <button on:click={save} data-test=save-button>save</button>
            {/if}
            {#if opened}
                <button on:click={close} data-test=close-button>v</button>
            {:else}
                <button on:click={open} data-test=open-button>&gt;</button>
            {/if}
        </span>
    </header>
    {#if opened}
        <div class='panel' transition:slide='{{ duration: 100 }}'>
            {#if editing}
                <textarea bind:value={task.description} placeholder='description' data-test=task-description-input></textarea>
            {:else}
                <p data-test=task-description>{@html (task.description || '')}</p>
            {/if}
        </div>
    {/if}
</section>