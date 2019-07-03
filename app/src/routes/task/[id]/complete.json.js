import * as taskRepo from '../../../data/task.repo'
import * as apiProxy from '../../../data/api.proxy'


export function put(req, res) {
	apiProxy.proxy(res, taskRepo.complete(req.params.id))
}