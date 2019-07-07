import fetch from 'node-fetch'
import { getToken } from '../auth/auth0'

const baseUrl = `http://${process.env.APPLICATION_HOST || 'localhost'}:${process.env.APPLICATION_PORT || '8080'}/api/v1`

export async function proxy(req, res, url) {
	const opts = {
		method: req.method,
		headers: req.headers,
	}
	if (Object.keys(req.body).length !== 0) {
		opts.body = JSON.stringify(req.body)
	}
	if (url === undefined) {
		url = req.url.replace('.json', '')
	}
	await checkAuth(opts)

	fetch(`${baseUrl}${url}`, opts).then(data => {
		res.writeHead(data.status, data.headers)
		data.text().then(text => {
			res.end(text)
		})
	}).catch(err => error(res, 500, err))
}

function error(res, code, msg) {
	res.writeHead(code, {
		'content-type': 'application/json'
	})
	res.end(JSON.stringify({
		error: msg
	}))
}

async function checkAuth(opts) {
	if (!opts.headers.authorization || opts.headers.authorization === 'Bearer undefined') {
		// Auth header doesn't exist from user session, use machine-to-machine Auth0 token for basic API auth
		const token = await getToken()
		opts.headers.authorization = `Bearer ${token}`
	}
}