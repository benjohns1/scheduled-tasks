import { getConfig } from '../../auth/auth0'

export async function get(_, res) {
  res.writeHead(200, {
		'Content-Type': 'application/json'
	})
  res.end(JSON.stringify(await getConfig()))
}