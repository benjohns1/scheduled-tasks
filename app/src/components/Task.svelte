<script>
    import { slide } from 'svelte/transition'
	import Button from "./Button.svelte"

    export let task = {}
    export let opened = false
    export let editing = undefined
    export let addTaskHandler = undefined
    export let completeTaskHandler = undefined

    if (!task.name) {
        task.name = ''
    }
    if (!task.description) {
        task.description = ''
    }
    
    function open() {
        opened = true
    }

    function close() {
        opened = false
    }

    function save() {
        if (addTaskHandler) {
            addTaskHandler(editing, task)
        }
    }

    function complete() {
        if (completeTaskHandler && task.id) {
            completeTaskHandler(task.id)
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
    .left {
        float: left;
        margin-right: 1rem;
    }
    header:after {
        content: "";
        clear: both;
        display: table;
    }
    textarea {
        width: 100%;
        height: 100%;
    }
    .card-text {
        padding-top: 0.5rem;
    }
</style>

<section class=card>
    <div class=card-body>
        <header>
            {#if !editing && !task.completedTime && completeTaskHandler}
                <span class=left>
                    <Button on:click={complete} test=complete-toggle style=outline-primary>done</Button>
                </span>
            {/if}
            <h3 data-test=task-name class='card-title form-inline'>
                {#if editing}
                    <input type=text bind:value={task.name} class=form-control placeholder='task name' data-test=task-name-input>
                {:else}
                    {(task.name || 'task')}
                {/if}
            </h3>
            <span class=right>
                {#if editing && addTaskHandler}
                    <Button on:click={save} test=save-button style=success>save</Button>
                {/if}
                {#if opened}
                    <Button on:click={close} test=close-button style=secondary>v</Button>
                {:else}
                    <Button on:click={open} test=open-button style=secondary>&gt;</Button>
                {/if}
            </span>
        </header>
        {#if opened}
            <div class=card-text transition:slide='{{ duration: 100 }}'>
                {#if editing}
                    <textarea bind:value={task.description} class=form-control placeholder='description' data-test=task-description-input></textarea>
                {:else}
                    <span data-test=task-description>{@html (task.description || '')}</span>
                {/if}
            </div>
        {/if}
    </div>
</section>