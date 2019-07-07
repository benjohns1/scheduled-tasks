import { proxy } from '../../../api/proxy'

export function put(req, res) {
	proxy(req, res, `/task/${req.params.id}/complete`)
}