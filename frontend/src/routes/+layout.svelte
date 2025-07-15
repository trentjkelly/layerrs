<script lang="ts">
    import { onMount } from 'svelte';
	import '../app.css';
    import AudioPlayer from '../components/AudioPlayer.svelte';
    import SideBar from '../components/SideBar.svelte';
    import { audio } from '../stores/player';
	import { jwt, refreshToken, isLoggedIn } from '../stores/auth';

	let { data, children } = $props();

	// Initialize audio component across the entire session
	onMount(async () => {
		audio.set(new Audio())

    	// Load tokens from cookies to session variables
		await loadCookies()
		console.log("jwt: " + $jwt)
		console.log("refreshToken: " + $refreshToken)

		const sessionKey = 'sessionStarted';

		await handleSessionStart()

		// Checking if this is the first page load for the session
		// if (!sessionStorage.getItem(sessionKey)) {
		// 	sessionStorage.setItem(sessionKey, 'true');
		// 	await handleSessionStart();
		// }
	});

	async function loadCookies() {
		console.log("loading cookies initially")
		if (data.newJWT) {
			console.log("Setting JWT")
        	jwt.set(data.newJWT)
			isLoggedIn.set(true)
    	}
    	if (data.newRefreshToken) {
			console.log("Setting refresh token")
        	refreshToken.set(data.newRefreshToken)
    	}
	}

	async function handleSessionStart() {
		console.log('Session started!');
		console.log($refreshToken)

		// Refresh token doesn't exist (already loaded from cookies in page.server.ts), then logout:
		if($refreshToken == "") {
			console.log("no refresh token")
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
			console.log("invalid refresh token, logging out")
			await deleteTokens()
			isLoggedIn.set(false)
		} // Refresh token was valid, so stay logged in
		else {
			console.log("refresh token is valid, so staying logged in")
			if (typeof newJWT !== 'string') {
				console.error("newJWT is not a string")
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
			console.error("Failed to set the JWT");
		}
	}

	async function refreshJWT() : Promise<(string | number)[]> {
		
		let newJWT = ""
		let status = 0

		try {
			// const backendURL = import.meta.env.VITE_BACKEND_URL;
			const res = await fetch(`https://layerrs.com/api/authentication/refresh`, {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({ $refreshToken })
			})
			status = await res.status
			const resJson = await res.json()
			newJWT = resJson.token
		} catch (error) {
			console.log("Request for refresh denied")
		}

		return [status, newJWT]
	}

	async function deleteTokens() {
		await fetch('/cookies', { method: 'DELETE' })
	}

</script>

<div class="h-screen w-screen flex flex-row bg-gray-900 text-white font-body">

	<SideBar></SideBar>

	{@render children()}

	<AudioPlayer></AudioPlayer>
    
</div>
