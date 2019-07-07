import fetch from 'node-fetch' 

if (process.env.NODE_ENV !== 'production') {
	require('dotenv').config({ path: '../.env'})
}

const { AUTH0_DOMAIN, AUTH0_WEBAPP_CLIENT_ID, AUTH0_API_IDENTIFIER } = process.env;

export const config = {
  "domain": AUTH0_DOMAIN,
  "clientId": AUTH0_WEBAPP_CLIENT_ID,
  "audience": AUTH0_API_IDENTIFIER,
}

const { AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET } = process.env;

let token = undefined

export async function getToken(forceNew = false) {
	if (!forceNew && token !== undefined) {
		return token
	}

	const authRequest = {
		"client_id": AUTH0_CLIENT_ID,
		"client_secret": AUTH0_CLIENT_SECRET,
		"audience": AUTH0_API_IDENTIFIER,
		"grant_type": "client_credentials"
	}

	await fetch(`https://dev-b1.auth0.com/oauth/token`, { method: "POST", body: JSON.stringify(authRequest), headers: { 'Content-Type': 'application/json' } }).then(async (data) => {
		await data.json().then(json => {
			token = json.access_token
		}).catch(err => console.error(err))
  }).catch(err => console.error(err))

  return token
}