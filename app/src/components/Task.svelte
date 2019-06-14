<script>
    import { slide } from 'svelte/transition';

    export let task = {};
    export let opened = false;
    export let editing = undefined;
    export let addTaskHandler = undefined;

    if (!task.name) {
        task.name = '';
    }
    if (!task.description) {
        task.description = '';
    }
    
    function open(event) {
        event.stopPropagation();
        opened = true;
    }

    function close(event) {
        event.stopPropagation();
        opened = false;
        editing = undefined;
    }

    function save(event) {
        event.stopPropagation();
        if (addTaskHandler) {
            addTaskHandler(editing, task);
        }
        editing = undefined;
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
    header span.right {
        float: right;
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

<section class='accordion'>
    <header>
        <h2 data-test='task-name'>
            {#if editing}
                <input type='text' bind:value={task.name} placeholder='task name' data-test='task-name-input'>
            {:else}
                {(task.name || 'task')}
            {/if}
        </h2>
        <span class='right'>
            {#if opened}
                {#if editing}
                    <button on:click={save} data-test='save-button'>save</button>
                {/if}
                <button on:click={close}>v</button>
            {:else}
                <button on:click={open} data-test='open-button'>></button>
            {/if}
        </span>
    </header>
    {#if opened}
        <div class='panel' transition:slide='{{ duration: 50 }}'>
            {#if editing}
                <textarea class='description' bind:value={task.description} placeholder='description' data-test='task-description-input'></textarea>
            {:else}
                <p class='description' data-test='task-description'>{@html (task.description || '')}</p>
            {/if}
        </div>
    {/if}
</section>