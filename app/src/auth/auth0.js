import fetch from 'node-fetch'

const {
	AUTH0_DOMAIN,
	AUTH0_WEBAPP_CLIENT_ID,
	AUTH0_API_IDENTIFIER,
	AUTH0_ANON_CLIENT_ID,
	AUTH0_ANON_CLIENT_SECRET,
	ENABLE_E2E_DEV_LOGIN,
	AUTH0_E2E_DEV_CLIENT_ID,
	AUTH0_E2E_DEV_CLIENT_SUBJECT,
	AUTH0_E2E_DEV_CLIENT_SECRET,
} = process.env;

export async function getConfig() {
	let cfg = {
		domain: AUTH0_DOMAIN,
		clientId: AUTH0_WEBAPP_CLIENT_ID,
		audience: AUTH0_API_IDENTIFIER,
	}
	if (ENABLE_E2E_DEV_LOGIN) {
		cfg.dev = {
			enabled: true,
			token: await getE2EDevToken(),
			subject: AUTH0_E2E_DEV_CLIENT_SUBJECT,
			displayname: "E2E Test User",
		}
	}
	return cfg
}

let anonymousToken = undefined
export async function getAnonymousToken(forceNew = false) {
	if (forceNew || !anonymousToken) {
		anonymousToken = await fetchToken(AUTH0_DOMAIN, AUTH0_ANON_CLIENT_ID, AUTH0_ANON_CLIENT_SECRET, AUTH0_API_IDENTIFIER)
	}
  return anonymousToken
}

async function fetchToken(domain, clientId, clientSecret, audience) {
	const authRequest = {
		"client_id": clientId,
		"client_secret": clientSecret,
		"audience": audience,
		"grant_type": "client_credentials"
	}

	return await fetch(`https://${domain}/oauth/token`, { method: "POST", body: JSON.stringify(authRequest), headers: { 'Content-Type': 'application/json' } })
		.then(data => data.json()
			.then(json => json.access_token)
			.catch(err => console.error(err)))
		.catch(err => console.error(err))
}

let devE2EToken = undefined
async function getE2EDevToken(forceNew = false) {
	if (!ENABLE_E2E_DEV_LOGIN) {
		console.error('attempting to call getE2EDevToken() when it is not enabled!')
		return null
	}
	if (forceNew || !devE2EToken) {
		devE2EToken = await fetchToken(AUTH0_DOMAIN, AUTH0_E2E_DEV_CLIENT_ID, AUTH0_E2E_DEV_CLIENT_SECRET, AUTH0_API_IDENTIFIER)
	}
	
	return devE2EToken
}