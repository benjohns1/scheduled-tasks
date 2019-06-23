import * as scheduleRepo from '../../../data/schedule.repo';
import * as apiProxy from '../../../data/api.proxy';

export function post(req, res) {
	apiProxy.proxy(res, scheduleRepo.addRecurringTask(req.params.id, req.body));
}