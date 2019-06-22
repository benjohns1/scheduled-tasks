import * as scheduleRepo from '../../../data/schedule.repo';
import * as apiProxy from '../../../data/api.proxy';

export function post(req, res) {
	if (req.params.command === 'task') {
		apiProxy.proxy(res, scheduleRepo.addRecurringTask(req.params.id, req.body));
	} else {
		apiProxy.error(res, 404, 'Not found');
	}
}

export function put(req, res) {
	switch (req.params.command) {
		case 'pause':
		case 'unpause':
			apiProxy.proxy(res, scheduleRepo.pause(req.params.id, req.params.command));
			break;
		default:
			apiProxy.error(res, 404, 'Not found');
	}
	
}