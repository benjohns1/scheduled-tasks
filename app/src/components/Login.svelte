<script context="module">
    import { onMount } from 'svelte'
    import createAuth0Client from '@auth0/auth0-spa-js'
	import { withJsonAndAuth } from "../api/default.headers"
    import Button from './Button.svelte'

    const loadConfig = async () => {
		return await fetch(`auth/config.json`).then(async r => {
			if (r.status === 200) {
				return { cfg: await r.json() }
			} else {
				throw {
					message: "Error retrieving auth config",
					r: await r.json()
				}
			}
		}).catch(error => {
            return { error }
        })
    }

    const authenticate = async (auth0) => {
        let isAuthenticated = await auth0.isAuthenticated()
        if (isAuthenticated) {
            return true
        }

        const query = window.location.search
        if (query.includes("code=") && query.includes("state=")) {
            await auth0.handleRedirectCallback()
            isAuthenticated = await auth0.isAuthenticated()
            window.history.replaceState({}, document.title, "/")
        }
        return isAuthenticated
    }
</script>

<script>
	import { loading } from './../loading-monitor'
    import * as sapper from '@sapper/app'
    const { session, page } = sapper.stores()
    if (!$session) {
        $session = {}
    }
    if (!$session.auth) {
        $session.auth = {}
    }
    
	const loaded = loading('login')
    
	let auth0 = undefined
    let errorMsg = undefined
    let config = {}
    let login = () => {}
    let logout = () => {}

    const loginHandler = () => login()
    const logoutHandler = () => logout()

    const getDevUser = () => {
        return {
            displayname: config.devDisplayname,
            sub: config.devSubject,
            iss: config.domain,
        }
    }

    const setCookies = (token, devLogin = false) => {
        const sessionExpireDays = 30
        const expires = new Date().getTime() + sessionExpireDays * 24 * 60 * 60 * 1000
        document.cookie = `token=${token}; expires=${expires}; path=/;`
        if (devLogin) {
            document.cookie = `devLogin=true; expires=${expires}; path=/;`
        } else {
            document.cookie = `devLogin=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/;`
        }
    }
    const clearCookies = () => {
        document.cookie = `token=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/;`
        document.cookie = `devLogin=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/;`
    }
    const devLogin = token => {
        $session.auth.devLogin = true
        $session.auth.token = token
        $session.auth.isAuthenticated = true
        $session.auth.user = getDevUser()
        setCookies(token, true)
        console.log('logged in as dev e2e test user')
        logout = () => {
            sessionLogout()
            sapper.goto('/')
            sessionAuth()
        }
        onUserLogin()
        sapper.goto($page.path)
    }

    const sessionLogin = async auth0 => {
        $session.auth.token = await auth0.getTokenSilently()
        $session.auth.user = await (async () => {
            const user = await auth0.getUser()
            user.displayname = user.nickname || user.name || user.email || 'New User'
            user.iss = config.domain
            return user
        })()
        setCookies($session.auth.token)
        logout = () => {
            auth0.logout()
            sessionLogout()
        }
        onUserLogin()
    }
    const sessionLogout = () => {
        $session.auth = {
            isAuthenticated: false
        }
        clearCookies()
    }

    const onUserLogin = async () => {
		return await fetch('auth/on-user-login.json', withJsonAndAuth($session, { method: 'POST', body: JSON.stringify($session.auth.user)})).then(async r => {
			if (r.status !== 204) {
				throw {
					message: "Error calling on-user-login hook",
					r: await r.text()
				}
			}
		}).catch(error => console.error(error))
    }

    const sessionAuth = async () => {
        // Load auth config from server
        const { cfg, error } = await loadConfig()
        if (error) {
            console.error(error)
            errorMsg = error.message || 'login error'
            return
        }
        config = cfg
        
        // Session logged-in as dev user - populate dev session data
        if ($session.auth.devLogin && $session.auth.token) {
            devLogin($session.auth.token)
            return
        }

        // Dev auto-login
        if ($page.query['dev-login'] && cfg.environment === "development" && cfg.token && process.browser) {
            devLogin(cfg.token)
            return
        }

        // Auth0 user login - prep user login methods Auth0
        auth0 = await createAuth0Client({
			domain: cfg.domain,
            client_id: cfg.clientId,
            audience: cfg.audience,
        })

        $session.auth.isAuthenticated = await authenticate(auth0)

        if ($session.auth.isAuthenticated) {
            // Get user data and token from Auth0 if user is currently logged-in
            await sessionLogin(auth0)
        }
        
        login = async () => {
            await auth0.loginWithRedirect({
                redirect_uri: window.location.origin
            })
            sessionLogin(auth0)
        }
    }

    onMount(async () => {
        await sessionAuth()
        loaded()
    })
</script>

<style>
    .login-text {
        vertical-align: bottom;
    }
    img.picture {
        border-radius: 50%;
    }
</style>

{#if $session && $session.auth && $session.auth.isAuthenticated}
    {#if errorMsg}
        <span class='text-danger login-text'>{errorMsg}</span>
    {/if}
    {#if $session.auth.user}
        <span class=login-text>{$session.auth.user.displayname}</span>
        {#if $session.auth.user.picture}
            <img class=picture src={$session.auth.user.picture} width=32 alt='user picture'/>
        {/if}
    {/if}
    <Button on:click={logoutHandler} style=outline-secondary classes=btn-sm>log out</Button>
{:else}
    <Button on:click={loginHandler} test=login-button style=outline-success classes=btn-sm disabled={!$session || !$session.auth || $session.auth.isAuthenticated === undefined}>log in | sign up</Button>
{/if}