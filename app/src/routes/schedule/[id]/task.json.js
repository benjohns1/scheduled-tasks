import * as taskRepo from '../../../data/schedule.repo';
import * as apiProxy from '../../../data/api.proxy';


export function post(req, res) {
	apiProxy.proxy(res, taskRepo.addRecurringTask(req.params.id, req.body));
}