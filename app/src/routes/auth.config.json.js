import { config } from '../auth/auth0'

export function get(_, res) {
  res.writeHead(200, {
		'Content-Type': 'application/json'
	})
  res.end(JSON.stringify(config))
}