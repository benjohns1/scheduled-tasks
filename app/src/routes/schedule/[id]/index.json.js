import * as scheduleRepo from '../../../data/schedule.repo'
import * as apiProxy from '../../../data/api.proxy'

export function del(req, res) {
  apiProxy.proxy(res, scheduleRepo.remove(req.params.id))
}