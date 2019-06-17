export function proxy(res, apiPromise) {
  apiPromise.then(data => {
		res.writeHead(data.status, data.headers);
		data.text().then(text => {
			res.end(text);
		});
	}).catch(err => {
		res.writeHead(500, {
      'Content-Type': 'application/json'
    });
		res.end(JSON.stringify({
			error: err
		}));
	});
}