import * as scheduleRepo from '../../data/schedule.repo'
import * as apiProxy from '../../data/api.proxy'

export function get(_, res) {
	apiProxy.proxy(res, scheduleRepo.getAll())
}

export function post(req, res) {
	apiProxy.proxy(res, scheduleRepo.add(req.body))
}