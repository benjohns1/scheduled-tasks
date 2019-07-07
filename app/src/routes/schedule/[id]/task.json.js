import { proxy } from '../../../api/proxy'

export function post(req, res) {
	proxy(req, res, `/schedule/${req.params.id}/task`)
}