<script>
	export let segment;
</script>

<style>
	nav {
		padding: 1rem 0;
	}
</style>

<nav>
	<ul class='nav nav-tabs'>
		<li class=nav-item><a class='nav-link {segment === undefined ? "active" : ""}' href='.' data-test='nav-home'>home</a></li>
		<li class=nav-item><a class='nav-link {segment === "task" ? "active" : ""}' href='task' data-test='nav-task'>tasks</a></li>
		<li class=nav-item><a class='nav-link {segment === "schedule" ? "active" : ""}' href='schedule' data-test='nav-schedule'>schedules</a></li>
	</ul>
</nav>