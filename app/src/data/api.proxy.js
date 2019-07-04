export function proxy(res, apiPromise) {
  apiPromise.then(data => {
		res.writeHead(data.status, data.headers)
		data.text().then(text => {
			res.end(text)
		})
	}).catch(err => error(res, 500, err))
}

export function error(res, code, msg) {
	res.writeHead(code, {
		'Content-Type': 'application/json'
	})
	res.end(JSON.stringify({
		error: msg
	}))
}

export async function getToken() {
	let jwt
	await fetch(`http://${process.env.APPLICATION_HOST || 'localhost'}:${process.env.APPLICATION_PORT || '8080'}/api/v1/auth/token`).then(data => {
		res.writeHead(data.status, data.headers)
		data.text().then(text => {
			jwt = text
			res.end(text)
		})
	}).catch(err => error(res, 500, err))
	return jwt
}