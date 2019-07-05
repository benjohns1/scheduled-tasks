import fetch from 'node-fetch'

if (process.env.NODE_ENV !== 'production') {
	require('dotenv').config({ path: '../.env'})
}
const { AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, AUTH0_API_IDENTIFIER } = process.env;

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
			token = {
				value: json.access_token,
				type: json.token_type,
			}
		}).catch(err => console.error(err))
  }).catch(err => console.error(err))

  return token
}

export async function withAuth(fetchOptions = { headers: {} }) {
  if (fetchOptions.headers === undefined) {
    fetchOptions.headers = {}
  }
  const token = await getToken()
  fetchOptions.headers['Authorization'] = `${token.type} ${token.value}`
  return fetchOptions
}