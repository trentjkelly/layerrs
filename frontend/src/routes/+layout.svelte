<script lang="ts">
	import '../app.css';
    import { onMount } from 'svelte';
    import AudioPlayer from '../components/AudioPlayer.svelte';
    import SideBar from '../components/SideBar.svelte';
    import { initializeAudio, audio } from '../stores/player';
	import { jwt, refreshToken, isLoggedIn } from '../stores/auth';
	import { handleEnvironment, urlBase } from '../stores/environment';
	import { logger } from '../lib/logger';
	
	let { data, children } = $props();

	function togglePlayPause() {
		if ($audio) {
			$audio.paused ? $audio.play() : $audio.pause()
		}
	}

	function rewindTrack(seconds: number) {
		if ($audio) {
			$audio.currentTime -= seconds
		}
	}

	function fastForwardTrack(seconds: number) {
		if ($audio) {
			$audio.currentTime += seconds
		}
	}
	
	// Initialize audio component across the entire session
	onMount(async () => {
		await handleEnvironment()
		initializeAudio()
		await loadCookies()
		await handleSessionStart()
		// handleHotkeys()
	});

	async function loadCookies() {
		if (data.newJWT) {
        	jwt.set(data.newJWT)
			isLoggedIn.set(true)
    	}
    	if (data.newRefreshToken) {
        	refreshToken.set(data.newRefreshToken)
    	}
	}

	// function handleHotkeys() {
	// 	const spaceDown = (e: KeyboardEvent) => {
	// 		if (e.key === ' ' || e.key === 'Space') {
	// 			e.preventDefault()
	// 			togglePlayPause()
	// 		}
	// 	}

	// 	const leftArrowDown = (e: KeyboardEvent) => {
	// 		if (e.key === 'ArrowLeft') {
	// 			e.preventDefault()
	// 			rewindTrack(5)
	// 		}
	// 	}

	// 	const rightArrowDown = (e: KeyboardEvent) => {
	// 		if (e.key === 'ArrowRight') {
	// 			e.preventDefault()
	// 			fastForwardTrack(5)
	// 		}
	// 	}

	// 	const letterKeyDown = (e: KeyboardEvent) => {
	// 		// Only handle single letter keys (a-z)
	// 		if (e.key.length === 1 && /^[a-z]$/i.test(e.key)) {
	// 			e.preventDefault()
	// 			jumpToSample(e.key.toLowerCase())
	// 		}
	// 	}

	// 	document.addEventListener('keydown', spaceDown)
	// 	document.addEventListener('keydown', leftArrowDown)
	// 	document.addEventListener('keydown', rightArrowDown)
	// 	document.addEventListener('keydown', letterKeyDown)
	// }

	// function jumpToSample(letter: string) {
	// 	if (!$audio || !$audio.duration) return
		
	// 	// Create a mapping for qwerty keyboard layout
	// 	const keyMap: { [key: string]: number } = {
	// 		'q': 0, 'w': 1, 'e': 2, 'r': 3, 't': 4, 'y': 5, 'u': 6, 'i': 7, 'o': 8, 'p': 9,
	// 		'a': 10, 's': 11, 'd': 12, 'f': 13, 'g': 14, 'h': 15, 'j': 16, 'k': 17, 'l': 18,
	// 		'z': 19, 'x': 20, 'c': 21, 'v': 22, 'b': 23, 'n': 24, 'm': 25
	// 	}
		
	// 	const letterIndex = keyMap[letter]
	// 	if (letterIndex === undefined) return
		
	// 	const position = letterIndex / 26
		
	// 	$audio.currentTime = $audio.duration * position
	// }

	async function handleSessionStart() {
		logger.debug('Session started!');
		logger.debug($refreshToken)

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
			logger.debug("invalid refresh token, logging out")
			await deleteTokens()
			isLoggedIn.set(false)
		} 
		// Refresh token was valid, so stay logged in
		else {
			logger.debug("refresh token is valid, so staying logged in")
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
				body: JSON.stringify({ $refreshToken })
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

<div class="h-screen w-screen flex flex-row bg-gray-900 text-white font-body">
	<SideBar></SideBar>
	{@render children()}
	<AudioPlayer></AudioPlayer>
</div>
