import * as taskRepo from '../../data/task.repo';

const contentType = {
	'Content-Type': 'application/json'
};

export function post(_, res) {
	taskRepo.clear().then(data => {
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