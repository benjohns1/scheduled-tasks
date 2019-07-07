import { proxy } from '../../api/proxy'

export function get(req, res) {
	proxy(req, res)
}

export function post(req, res) {
	proxy(req, res)
}