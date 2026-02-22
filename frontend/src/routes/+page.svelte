<script lang="ts">
    import { onMount } from "svelte";
    import TopHeader from "../components/TopHeader.svelte";
    import NewTrackCard from "../components/NewTrackCard.svelte";
    import UserMenu from "../components/UserMenu.svelte";
    import { isSidebarOpen } from "../stores/player";
    import { handleEnvironment, urlBase } from "../stores/environment";
    import { audio } from "../stores/player";

    // Each of the songs to be loaded in
    let artistId = 20 // Static for now

    /**
     * @type {any[]}
     */
    let trackIds = [];

    async function fetchData() {
        const response = await fetch(`${$urlBase}/api/recommendations/home`)
        const data = await response.json();
        trackIds = Object.keys(data).map(key => data[key])
    }

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

    function handleHotkeys() {
		const spaceDown = (e: KeyboardEvent) => {
			if (e.key === ' ' || e.key === 'Space') {
				e.preventDefault()
				togglePlayPause()
			}
		}

		const leftArrowDown = (e: KeyboardEvent) => {
			if (e.key === 'ArrowLeft') {
				e.preventDefault()
				rewindTrack(5)
			}
		}

		const rightArrowDown = (e: KeyboardEvent) => {
			if (e.key === 'ArrowRight') {
				e.preventDefault()
				fastForwardTrack(5)
			}
		}

		const letterKeyDown = (e: KeyboardEvent) => {
			// Only handle single letter keys (a-z)
			if (e.key.length === 1 && /^[a-z]$/i.test(e.key)) {
				e.preventDefault()
				jumpToSample(e.key.toLowerCase())
			}
		}

		document.addEventListener('keydown', spaceDown)
		document.addEventListener('keydown', leftArrowDown)
		document.addEventListener('keydown', rightArrowDown)
		document.addEventListener('keydown', letterKeyDown)
	}

	function jumpToSample(letter: string) {
		if (!$audio || !$audio.duration) return
		
		// Create a mapping for qwerty keyboard layout
		const keyMap: { [key: string]: number } = {
			'q': 0, 'w': 1, 'e': 2, 'r': 3, 't': 4, 'y': 5, 'u': 6, 'i': 7, 'o': 8, 'p': 9,
			'a': 10, 's': 11, 'd': 12, 'f': 13, 'g': 14, 'h': 15, 'j': 16, 'k': 17, 'l': 18,
			'z': 19, 'x': 20, 'c': 21, 'v': 22, 'b': 23, 'n': 24, 'm': 25
		}
		
		const letterIndex = keyMap[letter]
		if (letterIndex === undefined) return
		
		const position = letterIndex / 26
		
		$audio.currentTime = $audio.duration * position
	}

    onMount(async () => {
        await handleEnvironment()
        await fetchData()
        // handleHotkeys()
    })

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }

</script>

<main class={`transition-all duration-300 h-auto w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>

    <TopHeader pageName="Home" pageIcon="home.png"></TopHeader>

    <!-- User Menu -->
    <UserMenu />

    <!-- Where the songs go -->
    <section class="w-full flex flex-wrap justify-around pb-24">
        <!-- {#each trackIds as id}
            <NewTrackCard trackId={id}></NewTrackCard>
        {/each} -->
        <NewTrackCard trackId={artistId}></NewTrackCard>
    </section>
</main>
