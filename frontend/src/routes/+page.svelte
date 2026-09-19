<script lang="ts">
    import { page } from "$app/stores";
    import TopHeader from "../components/TopHeader.svelte";
    import { urlBase } from "../stores/environment";
    import { audio } from "../stores/player";
    import { authInitialized } from "../stores/auth";
    import { fetchWithAuth } from "../modules/lib/fetch";
    import type { TrackInfo, TrackPagesBulkResponse } from "../models/types";
    import { getTrackPagesBulk } from "../modules/requests/page-requests";
    import { getTrackUsesBulk, getTrackInfoBatch } from "../modules/requests/track-requests";
    import ProjectTrackCard from "../components/ProjectTrackCard.svelte";

    let tracks: TrackInfo[] = $state([]);
    let trackPages: Map<number, TrackPagesBulkResponse> = $state(new Map());
    let trackUses: Map<number, TrackInfo[]> = $state(new Map());
    let searchTerm = $derived(($page.url.searchParams.get('search') ?? '').trim());
    let isLoading = $state(false);
    let errorMessage = $state('');
    let activeRequestId = 0;

    async function fetchData(query: string, requestId: number) {
        const endpoint = query
            ? `${$urlBase}/api/search?query=${encodeURIComponent(query)}`
            : `${$urlBase}/api/recommendations/home`;
        const response = await fetchWithAuth(endpoint);
        if (!response.ok) {
            throw new Error(`Failed to load tracks (${response.status})`);
        }
        const data = await response.json() as TrackInfo[];
        if (requestId !== activeRequestId) return;

        tracks = data;
        await fetchTrackPages(data, requestId);
        await fetchTrackUses(data, requestId);
    }

    async function fetchTrackPages(currentTracks: TrackInfo[], requestId: number) {
        if (currentTracks.length === 0) {
            return;
        }

        const trackIds = currentTracks.map(t => t.id);
        const response = await getTrackPagesBulk($urlBase, trackIds);
        if (response && requestId === activeRequestId) {
            const pagesMap = new Map<number, TrackPagesBulkResponse>();
            for (const [trackId, data] of Object.entries(response)) {
                pagesMap.set(Number(trackId), data);
            }
            trackPages = pagesMap;
        }
    }

    async function fetchTrackUses(currentTracks: TrackInfo[], requestId: number) {
        if (currentTracks.length === 0) {
            return;
        }

        const trackIds = currentTracks.map(t => t.id);
        const usesResponse = await getTrackUsesBulk($urlBase, trackIds);
        if (!usesResponse || requestId !== activeRequestId) {
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
            if (requestId !== activeRequestId) return;
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
        if (requestId === activeRequestId) {
            trackUses = usesMap;
        }
    }

    $effect(() => {
        if ($authInitialized) {
            const requestId = ++activeRequestId;
            tracks = [];
            trackPages = new Map();
            trackUses = new Map();
            isLoading = true;
            errorMessage = '';

            fetchData(searchTerm, requestId)
                .catch((error) => {
                    if (requestId === activeRequestId) {
                        errorMessage = error instanceof Error ? error.message : 'Unable to load tracks.';
                    }
                })
                .finally(() => {
                    if (requestId === activeRequestId) {
                        isLoading = false;
                    }
                });
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

</script>

<main class="h-auto w-full bg-zinc-950">

    <TopHeader pageName="Home" pageIcon="home.png"></TopHeader>

    <!-- Where the songs go -->
    <section class="w-full pb-24">
        <div class="w-full px-4 sm:px-6 lg:px-8">
			{#if isLoading}
				<div class="border border-zinc-800 bg-zinc-900 px-6 py-12 text-center text-zinc-400">
					Loading tracks…
				</div>
			{:else if errorMessage}
				<div class="border border-red-900 bg-zinc-900 px-6 py-12 text-center">
					<p class="text-lg font-semibold text-white">Unable to load tracks</p>
					<p class="mt-2 text-sm text-zinc-400">{errorMessage}</p>
				</div>
			{:else}
			{#each tracks as track, index}
				{@const pagesData = trackPages.get(track.id)}
			{@const usesData = trackUses.get(track.id)}
				 <ProjectTrackCard
					{track}
					trackNumber={index + 1}
					pages={pagesData?.pages ?? []}
					pageCount={pagesData?.pageCount ?? 0}
					uses={usesData ?? []}
				></ProjectTrackCard>
			{/each}
			{#if searchTerm && tracks.length === 0}
				<div class="border border-zinc-800 bg-zinc-900 px-6 py-12 text-center">
					<p class="text-lg font-semibold text-white">No tracks found</p>
					<p class="mt-2 text-sm text-zinc-400">Try a different track or artist name.</p>
				</div>
			{/if}
			{/if}
		</div>
		
    </section>
</main>
