<script>
    export let test = undefined
    export let style = 'primary'
    export let classes = ''
    export let disabled = false

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
    <button type=button {disabled} on:click data-test={test} class='btn btn-{style}{classStr}'><slot/></button>
{:else}
    <button type=button {disabled} on:click class='btn btn-{style}{classStr}'><slot/></button>
{/if}