<script lang="ts">
    import { onMount } from 'svelte';
	import '../app.css';
    import AudioPlayer from '../components/AudioPlayer.svelte';
    import SideBar from '../components/SideBar.svelte';
    import { audio } from '../stores/player';
	import { refreshToken, isLoggedIn } from '../stores/auth';

	let { children } = $props();

	// Initialize audio component across the entire session
	onMount(async () => {
		audio.set(new Audio())

		const sessionKey = 'sessionStarted';

		// Checking if this is the first page load for the session
		if (!sessionStorage.getItem(sessionKey)) {
			sessionStorage.setItem(sessionKey, 'true');
			await handleSessionStart();
		}
	});

	async function handleSessionStart() {
		console.log('Session started!');

		// Refresh JWT if it's not expired
		await refreshJWT($refreshToken)
	}

	async function refreshJWT(refreshToken: string) {

		if(!refreshToken) {
			isLoggedIn.set(false)
			return
		}
		
		let newJWT = ""
		let newRefreshToken = ""

		console.log(refreshToken)

		try {
			const res = await fetch(`http://localhost:8080/api/authentication/refresh`, {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({ refreshToken })
			})

			const resJson = await res.json()
			newJWT = resJson.token
			newRefreshToken = resJson.refreshToken
		} catch (error) {
			console.log("Request for refresh denied")
		}

		console.log(newJWT)
		console.log(newRefreshToken)

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

</script>

<div class="h-screen w-screen flex flex-row bg-gray-900 text-white">

	<SideBar></SideBar>

	{@render children()}

	<AudioPlayer></AudioPlayer>
    
</div>
