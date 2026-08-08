<script lang="ts">
    import { onMount } from "svelte";
    import TopHeader from "../components/TopHeader.svelte";
	import UserMenu from "../components/UserMenu.svelte";
    import { isSidebarOpen } from "../stores/player";
    import { urlBase } from "../stores/environment";
    import { usernameStore, emailStore, portraitUrlStore } from "../stores/profile";
    import { audio } from "../stores/player";
    import { authInitialized } from "../stores/auth";
    import { fetchWithAuth } from "../modules/lib/fetch";
    import type { TrackInfo, TrackPagesBulkResponse } from "../models/types";
    import { getTrackPagesBulk } from "../modules/requests/page-requests";
    import { getTrackUsesBulk, getTrackInfoBatch } from "../modules/requests/track-requests";
    import FinalTrackCard from "../components/FinalTrackCard.svelte";

    let tracks: TrackInfo[] = $state([]);
    let trackPages: Map<number, TrackPagesBulkResponse> = $state(new Map());
    let trackUses: Map<number, TrackInfo[]> = $state(new Map());

    async function fetchData() {
        const response = await fetchWithAuth(`${$urlBase}/api/recommendations/home`)
        const data = await response.json();
        tracks = data as TrackInfo[]
        await fetchTrackPages();
        await fetchTrackUses();
    }

    async function fetchTrackPages() {
        if (tracks.length === 0) {
            trackPages = new Map();
            return;
        }

        const trackIds = tracks.map(t => t.id);
        const response = await getTrackPagesBulk($urlBase, trackIds);
        if (response) {
            const pagesMap = new Map<number, TrackPagesBulkResponse>();
            for (const [trackId, data] of Object.entries(response)) {
                pagesMap.set(Number(trackId), data);
            }
            trackPages = pagesMap;
        }
    }

    async function fetchTrackUses() {
        if (tracks.length === 0) {
            trackUses = new Map();
            return;
        }

        const trackIds = tracks.map(t => t.id);
        const usesResponse = await getTrackUsesBulk($urlBase, trackIds);
        if (!usesResponse) {
            trackUses = new Map();
            return;
        }

        const allUseIds = new Set<number>();
        for (const useIds of Object.values(usesResponse)) {
            for (const useId of useIds) {
                allUseIds.add(useId);
            }
        }

        let useTracksMap = new Map<number, TrackInfo>();
        if (allUseIds.size > 0) {
            const useTracks = await getTrackInfoBatch($urlBase, Array.from(allUseIds));
            if (useTracks) {
                for (const useTrack of useTracks) {
                    useTracksMap.set(useTrack.id, useTrack);
                }
            }
        }

        const usesMap = new Map<number, TrackInfo[]>();
        for (const [trackId, useIds] of Object.entries(usesResponse)) {
            const uses: TrackInfo[] = [];
            for (const useId of useIds) {
                const useTrack = useTracksMap.get(useId);
                if (useTrack) {
                    uses.push(useTrack);
                }
            }
            usesMap.set(Number(trackId), uses);
        }
        trackUses = usesMap;
    }

    $effect(() => {
        if ($authInitialized) {
            fetchData();
        }
    });

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
        await fetchData()
        // handleHotkeys()
    })

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }

</script>

<main class={`transition-all duration-300 h-auto w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-olive-400`}>

    <TopHeader pageName="Home" pageIcon="home.png"></TopHeader>

    <!-- User Menu -->
    <UserMenu username={$usernameStore} email={$emailStore} portraitUrl={$portraitUrlStore} />

    <!-- Where the songs go -->
    <section class="w-full flex flex-wrap justify-around pb-24">
        <div class="h-full w-1/2">
			{#each tracks as track}
				<!-- <NewTrackCard {track}></NewTrackCard> -->
				{@const pagesData = trackPages.get(track.id)}
			{@const usesData = trackUses.get(track.id)}
				 <FinalTrackCard
					{track}
					pages={pagesData?.pages ?? []}
					pageCount={pagesData?.pageCount ?? 0}
					uses={usesData ?? []}
				></FinalTrackCard>
			{/each}
		</div>
		
		<!-- <FinalTrackCard track={tracks[0]}></FinalTrackCard> -->
    </section>
</main>
