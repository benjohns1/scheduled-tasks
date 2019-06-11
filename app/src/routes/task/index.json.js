import * as taskRepo from '../../data/task.repo';

// Pass-through tasks from API service to Sapper Node API
export async function get(_, res) {
	const data = taskRepo.getAll();

	res.writeHead(200, {
		'Content-Type': 'application/json'
	});
	
	res.end(await data);
}