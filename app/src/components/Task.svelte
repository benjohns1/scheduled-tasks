<script>
    export let task;
    export let open = false;

    import { slide } from 'svelte/transition';

    function toggle() {
        open = !open
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
    header .toggler {
        cursor: pointer;
    }
    header .toggler span {
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
</style>

<section class="accordion">
    <header>
        <div class="toggler" on:click={toggle}>
            <h2>{task.name}</h2>
            <span>{#if open}v{:else}>{/if}</span>
        </div>
    </header>
    {#if open}
        <div class="panel" transition:slide="{{ duration: 50 }}">
            <p>{@html task.description}</p>
        </div>
    {/if}
</section>