<script context="module">
    import { onMount } from 'svelte'
    import createAuth0Client from '@auth0/auth0-spa-js'
    import Button from './Button.svelte'

    const loadConfig = async () => {
		return await fetch(`auth.config.json`).then(async r => {
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
    import { stores } from '@sapper/app'
    const { session } = stores()
    if (!$session) {
        $session = {
            isAuthenticated: undefined,
        }
    }
    
	let auth0 = undefined
    let errorMsg = undefined

    onMount(async () => {
        const { cfg, error } = await loadConfig()
        if (error) {
            console.error(error)
            errorMsg = error.message || 'login error'
            return
        }

        auth0 = await createAuth0Client({
			domain: cfg.domain,
            client_id: cfg.clientId,
            audience: cfg.audience,
        })

        $session.isAuthenticated = await authenticate(auth0)

        if ($session.isAuthenticated) {
            $session.user = await getUser()
            $session.token = await getToken()
            
            const expireDays = 30
            document.cookie = `token=${$session.token}; expires=${(new Date()).getTime() + (expireDays*24*60*60*1000)}; path=/;`
        }
    })

	const login = async () => {
		await auth0.loginWithRedirect({
            redirect_uri: window.location.origin
        })
	}

	const logout = () => {
		auth0.logout({
            returnTo: window.location.origin
        })
        document.cookie = `token=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/;`
    }

    const getUser = async () => {
        const user = await auth0.getUser()
        user.displayname = user.nickname || user.name || user.email || 'New User'
        return user
    }

    const getToken = async () => {
        return await auth0.getTokenSilently()
    }

</script>

<style>
    .login-text {
        vertical-align: bottom;
    }
    img.picture {
        border-radius: 50%;
    }
</style>

{#if $session && $session.isAuthenticated}
    {#if errorMsg}
        <span class='text-danger login-text'>{errorMsg}</span>
    {/if}
    {#if $session.user}
        <span class=login-text>{$session.user.displayname}</span>
        {#if $session.user.picture}
            <img class=picture src={$session.user.picture} width=32 alt='user picture'/>
        {/if}
    {/if}
    <Button on:click={logout} style=outline-secondary classes=btn-sm>log out</Button>
{:else}
    <Button on:click={login} style=outline-success classes=btn-sm disabled={!$session || $session.isAuthenticated === undefined}>log in</Button>
{/if}