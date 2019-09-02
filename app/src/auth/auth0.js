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

const E2E_DEV_LOGIN_LOCAL_EXPIRE = 600

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
	if (!domain || !clientId || !clientSecret || !audience) {
		console.error('Error: at least 1 auth config params is empty, could not retrieve token')
		return ''
	}

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
let devE2ETokenCacheTime = 0
async function getE2EDevToken(forceNew = false) {
	if (!ENABLE_E2E_DEV_LOGIN) {
		console.error('attempting to call getE2EDevToken() when it is not enabled!')
		return null
	}
	const nowSeconds = Math.floor(Date.now() / 1000)
	if (devE2ETokenCacheTime + E2E_DEV_LOGIN_LOCAL_EXPIRE < nowSeconds) {
		// Local expiration time is up, refresh the token
		devE2EToken = undefined
	}
	if (forceNew || !devE2EToken) {
		console.log('fetching new E2E Dev Token')
		devE2EToken = await fetchToken(AUTH0_DOMAIN, AUTH0_E2E_DEV_CLIENT_ID, AUTH0_E2E_DEV_CLIENT_SECRET, AUTH0_API_IDENTIFIER)
		devE2ETokenCacheTime = nowSeconds
	}

	return devE2EToken
}
