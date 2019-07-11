import { proxy } from '../../../api/proxy'

export function del(req, res) {
	proxy(req, res, { url: `/schedule/${req.params.id}` })
}