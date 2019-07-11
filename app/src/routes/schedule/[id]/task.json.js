import { proxy } from '../../../api/proxy'

export function post(req, res) {
	proxy(req, res, { url: `/schedule/${req.params.id}/task` })
}