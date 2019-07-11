import fetch from 'node-fetch'
import { getAnonymousToken } from '../auth/auth0'

const baseUrl = `http://${process.env.APPLICATION_HOST || 'localhost'}:${process.env.APPLICATION_PORT || '8080'}/api/v1`

export async function proxy(req, res, { method = undefined, url = undefined, body = undefined } = {}) {
	if (method === undefined) {
		method = req.method
	}
	if (body === undefined) {
		body = req.body
	}
	if (url === undefined) {
		url = req.url.replace('.json', '')
	}

	const opts = {
		method: method,
		headers: req.headers,
	}
	if (Object.keys(req.body).length !== 0) {
		opts.body = JSON.stringify(req.body)
	}
	await checkAuth(opts)

	fetch(`${baseUrl}${url}`, opts).then(data => {
		res.writeHead(data.status, data.headers)
		data.text().then(text => {
			res.end(text)
		})
	}).catch(err => error(res, 500, err))
}

export function error(res, code, msg) {
	res.writeHead(code, {
		'content-type': 'application/json'
	})
	res.end(JSON.stringify({
		error: msg
	}))
}

async function checkAuth(opts) {
	if (!opts.headers.authorization || opts.headers.authorization === 'Bearer undefined') {
		// Auth header doesn't exist from user session, use machine-to-machine Auth0 token for anonymous user
		const anonymousToken = await getAnonymousToken()
		opts.headers.authorization = `Bearer ${anonymousToken}`
	}
}