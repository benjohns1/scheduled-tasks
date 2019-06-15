import * as taskRepo from '../../data/task.repo';

const contentType = {
	'Content-Type': 'application/json'
};

// Pass-through tasks from API service to Sapper Node API
export function get(_, res) {
	taskRepo.getAll().then(data => {
		res.writeHead(data.status, data.headers);
		data.text().then(text => {
			res.end(text);
		});
	}).catch(err => {
		res.writeHead(500, contentType);
		res.end(JSON.stringify({
			error: err
		}));
	});
}

export function post(req, res) {
	taskRepo.add(req.body).then(data => {
		res.writeHead(data.status, data.headers);
		data.text().then(text => {
			res.end(text);
		});
	}).catch(err => {
		res.writeHead(500, contentType);
		res.end(JSON.stringify({
			error: err
		}));
	});
}