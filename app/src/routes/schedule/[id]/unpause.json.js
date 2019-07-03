import * as scheduleRepo from '../../../data/schedule.repo'
import * as apiProxy from '../../../data/api.proxy'

export function put(req, res) {
	apiProxy.proxy(res, scheduleRepo.unpause(req.params.id))
}