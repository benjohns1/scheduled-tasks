import { proxy } from '../../../api/proxy'

export function del(req, res) {
	proxy(req, res, `/schedule/${req.params.id}`)
}