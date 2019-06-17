import * as taskRepo from '../../data/schedule.repo';
import * as apiProxy from '../../data/api.proxy';

export function get(_, res) {
	apiProxy.proxy(res, taskRepo.getAll());
}

export function post(req, res) {
	apiProxy.proxy(res, taskRepo.add(req.body));
}