<script>
    export let test = undefined
    export let style = 'primary'
    export let classes = ''

    $: classStr = classes === '' ? '' : ` ${classes}`
</script>

<style>
    .right {
        float: right;
    }
    .left {
        float: left;
    }
</style>

{#if test !== undefined}
    <button on:click data-test={test} class='btn btn-{style}{classStr}'><slot/></button>
{:else}
    <button on:click class='btn btn-{style}'><slot/></button>
{/if}