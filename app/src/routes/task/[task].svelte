<script context="module">
    import Task from "../../components/Task.svelte"

    export async function preload({ params, query }) {
        const res = await this.fetch(`task/${params.task}.json`);
        const data = await res.json();

        if (res.status === 200) {
            return { task: data };
        } else {
			this.error(res.status, data.message);
        }
    }
</script>

<script>
    export let task;
</script>

<svelte:head>
    <title>Scheduled Tasks - {task.name}</title>
</svelte:head>

<Task {task} open={true} />