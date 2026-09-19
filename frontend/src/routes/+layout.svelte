<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import AudioPlayer from '../components/AudioPlayer.svelte';
	import AppHeader from '../components/AppHeader.svelte';
	import { initializeAudio, audio } from '../stores/player';
	import { jwt, refreshToken, isLoggedIn, authInitialized } from '../stores/auth';
	import { urlBase, stripePublishableKey } from '../stores/environment';
	import { loadProfile } from '../stores/profile';
	import { logger } from '../modules/lib/logger';
	import { loadStripe } from '../modules/lib/stripe';

	let { data, children } = $props();
	let searchActive = $state(false);

	// Initialize audio component across the entire session
	onMount(async () => {
		initializeAudio();
		await loadCookies();
		if ($page.url.pathname === '/login/callback') {
			authInitialized.set(true);
			return;
		}
		await handleSessionStart();
		if ($isLoggedIn) {
			await loadProfile();
		}
		if ($stripePublishableKey) {
			await loadStripe($stripePublishableKey);
		}
		authInitialized.set(true);
	});

	$effect(() => {
		if ($authInitialized && !$isLoggedIn && !isPublicRoute($page.url.pathname)) {
			goto('/login', { replaceState: true });
		}
	});

	function isPublicRoute(pathname: string) {
		return pathname === '/' || pathname === '/login' || pathname.startsWith('/login/');
	}

	async function loadCookies() {
		if (data.newJWT) {
        	jwt.set(data.newJWT)
			isLoggedIn.set(true)
    	}
    	if (data.newRefreshToken) {
        	refreshToken.set(data.newRefreshToken)
    	}
	}

	async function handleSessionStart() {
		logger.debug('Session started!');

		// Refresh token doesn't exist (already loaded from cookies in page.server.ts), then logout:
		if($refreshToken == "") {
			await deleteTokens()
			isLoggedIn.set(false)
			return
		}

		// Refresh JWT w/ Refresh token (get back status code, potentially jwt)
		const values = await refreshJWT()
		const status = values[0]
		const newJWT = values[1]

		// Refresh token was invalid, so log out
		if (status == 401 || status == 500) {
			await deleteTokens()
			isLoggedIn.set(false)
		} 
		// Refresh token was valid, so stay logged in
		else {
			if (typeof newJWT !== 'string') {
				logger.error("newJWT is not a string")
			} else {
				await writeTokensToCookies(newJWT, $refreshToken)
				isLoggedIn.set(true)
			}
		}
	}

	async function writeTokensToCookies(newJWT : string, newRefreshToken : string) {
		try {
			const res2  = await fetch('/cookies', { 
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({ 
					jwtToken: newJWT,
					refreshToken: newRefreshToken
				})
			})
		} catch (error) {
			logger.error("Failed to set the JWT");
		}
	}

	async function refreshJWT() : Promise<(string | number)[]> {
		let newJWT = ""
		let status = 0

		try {
			const res = await fetch(`${$urlBase}/api/authentication/refresh`, {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({ refreshToken: $refreshToken })
			})
			status = await res.status
			const resJson = await res.json()
			newJWT = resJson.token
		} catch (error) {
			logger.error("Request for refresh denied")
		}

		return [status, newJWT]
	}

	async function deleteTokens() {
		await fetch('/cookies', { method: 'DELETE' })
	}

</script>

<div class="flex h-dvh w-screen flex-col overflow-hidden bg-zinc-950 font-body text-white">
	<AppHeader bind:searchActive />
	{#if $authInitialized && (isPublicRoute($page.url.pathname) || $isLoggedIn)}
		<div class={`app-content min-h-0 flex-1 overflow-y-auto overflow-x-hidden transition-[filter] duration-150 ${searchActive ? 'blur-[1.5px]' : ''}`}>
			{@render children()}
		</div>
	{/if}
</div>
