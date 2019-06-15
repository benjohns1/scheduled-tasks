import * as taskRepo from '../../../data/task.repo';

const contentType = {
	'Content-Type': 'application/json'
};

export function put(req, res) {
	taskRepo.complete(req.params.id).then(data => {
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