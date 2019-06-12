import * as taskRepo from '../../data/task.repo';

// Pass-through tasks from API service to Sapper Node API
export async function get(_, res) {

	const taskPromise = taskRepo.getAll();
	const contentType = {
		'Content-Type': 'application/json'
	};

	taskPromise.then(data => {
		res.writeHead(200, contentType);
		res.end(data);
	}).catch(err => {
		res.writeHead(500, contentType);
		res.end(JSON.stringify({
			error: err
		}));
	});
}