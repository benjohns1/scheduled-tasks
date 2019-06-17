import * as taskRepo from '../../data/task.repo';
import * as apiProxy from '../../data/api.proxy';

export function post(_, res) {
	apiProxy.proxy(res, taskRepo.clear());
}