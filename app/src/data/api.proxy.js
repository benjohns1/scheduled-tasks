export function proxy(res, apiPromise, useAuth = true) {
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