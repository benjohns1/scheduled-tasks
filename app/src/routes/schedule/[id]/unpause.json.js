import { proxy } from '../../../api/proxy'

export function put(req, res) {
	proxy(req, res, { url: `/schedule/${req.params.id}/unpause` })
}