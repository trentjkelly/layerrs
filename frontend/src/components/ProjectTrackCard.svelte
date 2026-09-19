<script lang="ts">
    import { onDestroy, onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase } from '../stores/environment';
    import type { PageWithFollowerCount, TrackInfo } from '../models/types';
    import { audio, currentTrack, currentTrackId, currentTime as globalCurrentTime, isPlaying } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import { trackPlayProgress } from '../modules/lib/play-tracking';
    import AddToPageButton from './AddToPageButton.svelte';
    import LikeButton from './LikeButton.svelte';
    import PlayPauseButton from './PlayPauseButton.svelte';
    import SourceStemLane from './SourceStemLane.svelte';
    import WaveformBar from './WaveformBar.svelte';

    let {
        track,
        trackNumber = 0,
        pages = [],
        pageCount = 0,
        uses = []
    }: {
        track: TrackInfo;
        trackNumber?: number;
        pages?: PageWithFollowerCount[];
        pageCount?: number;
        uses?: TrackInfo[];
    } = $props();

    let urlBase = $state('');
    let newAudioURL = $state('');
    let isTrackLiked = $state(false);
    let likeCount = $state(0);
    let isStemsExpanded = $state(false);
    let showSharePopup = $state(false);
    let copiedToClipboard = $state(false);
    let shareArea = $state<HTMLElement | null>(null);

    let waveformWidth = $state(0);
    let visibleBars = $state<number[]>([]);
    let timePerBar = $state(0);
    let cursorTime = $state(0);
    let isCursorHovered = $state(false);
    let currentSongPercentage = $state(0);

    $effect(() => {
        isTrackLiked = track.isLiked;
        likeCount = track.likes;
    });

    $effect(() => {
        updateWaveformBars();
    });

    $effect(() => {
        if (!$audio) return;

        const updateTime = () => {
            if ($currentTrackId === track.id) {
                const currentTime = $audio.currentTime;
                globalCurrentTime.set(currentTime);
                currentSongPercentage = currentTime / track.duration * 100;
            }
        };
        const handleEnded = () => {
            if ($currentTrackId === track.id) currentSongPercentage = 0;
        };

        $audio.addEventListener('timeupdate', updateTime);
        $audio.addEventListener('ended', handleEnded);
        return () => {
            $audio.removeEventListener('timeupdate', updateTime);
            $audio.removeEventListener('ended', handleEnded);
        };
    });

    $effect(() => {
        if ($currentTrackId === track.id && urlBase) return trackPlayProgress(track.id, urlBase);
    });

    $effect(() => {
        if (!showSharePopup) return;
        const closeOnOutsideClick = (event: MouseEvent) => {
            if (!shareArea?.contains(event.target as Node)) {
                showSharePopup = false;
                copiedToClipboard = false;
            }
        };
        document.addEventListener('click', closeOnOutsideClick);
        return () => document.removeEventListener('click', closeOnOutsideClick);
    });

    onMount(() => {
        urlBase = getUrlBase();
    });

    onDestroy(() => {
        if (newAudioURL) URL.revokeObjectURL(newAudioURL);
    });

    function updateWaveformBars() {
        const barStride = 4;
        const numBars = Math.floor(waveformWidth / barStride);
        if (numBars <= 0 || track.waveformData.length === 0) {
            visibleBars = [];
            timePerBar = 0;
            return;
        }
        visibleBars = Array.from({ length: numBars }, (_, index) => {
            const waveformIndex = Math.min(track.waveformData.length - 1, Math.floor(track.waveformData.length * index / numBars));
            return track.waveformData[waveformIndex];
        });
        timePerBar = track.duration / numBars;
    }

    function handleMouseMove(event: MouseEvent) {
        const bounds = (event.currentTarget as HTMLElement).getBoundingClientRect();
        const percentage = Math.max(0, Math.min(1, (event.clientX - bounds.left) / waveformWidth));
        cursorTime = percentage * track.duration;
    }

    async function playOrSeek() {
        if (!$audio) return;
        if ($currentTrackId === track.id) {
            isPlaying.set(true);
            $audio.currentTime = cursorTime;
            await $audio.play();
            return;
        }

        isPlaying.set(true);
        newAudioURL = await getAudio(urlBase, String(track.id));
        currentTrack.set(newAudioURL);
        currentTrackId.set(track.id);
        $audio.src = newAudioURL;
        try {
            await $audio.play();
        } catch (error) {
            logger.error(`Failed to play project: ${error}`);
        }
    }

    function toggleSharePopup(event: MouseEvent) {
        event.stopPropagation();
        showSharePopup = !showSharePopup;
        copiedToClipboard = false;
    }

    function copyLink(event: MouseEvent) {
        event.stopPropagation();
        navigator.clipboard.writeText(`${window.location.origin}/track/${track.id}`);
        copiedToClipboard = true;
    }

    function toggleStems(event: MouseEvent) {
        event.stopPropagation();
        isStemsExpanded = !isStemsExpanded;
    }

    function navigateLayerr() {
        goto(`/layerrs/${track.id}`, { state: { track: $state.snapshot(track) } });
    }

    function formatDuration(seconds: number): string {
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    }
</script>

<article class="group w-full overflow-visible border-b border-zinc-800 bg-zinc-950" aria-label="Project {track.description}">
    <div class="relative overflow-visible border border-zinc-700 bg-primary shadow-[inset_0_1px_rgba(255,255,255,0.12)]">
        <header class="relative z-20 flex min-h-8 items-center gap-2 border-b border-black/25 bg-black/15 px-2 text-xs text-violet-50">
            <span class="shrink-0 border-r border-black/20 pr-2 text-[10px] text-violet-100/80">
                {String(trackNumber || track.id).padStart(2, '0')}
            </span>
            <div class="min-w-0 flex flex-1 items-center gap-1 overflow-x-auto whitespace-nowrap [scrollbar-width:none]">
                {#each uses as sourceTrack, index (sourceTrack.id)}
                    <a href={`/track/${sourceTrack.id}`} class="shrink-0 font-medium hover:underline focus-visible:outline focus-visible:outline-1 focus-visible:outline-white">{sourceTrack.description}</a>
                    <span class="text-violet-100/70" aria-hidden="true">→</span>
                {/each}
                <a href={`/track/${track.id}`} class="shrink-0 font-semibold hover:underline focus-visible:outline focus-visible:outline-1 focus-visible:outline-white">{track.description}</a>
            </div>
            <a href={`/profile/${track.artistId}`} class="hidden max-w-36 shrink-0 truncate border-l border-black/20 pl-2 text-[10px] text-violet-100/85 hover:underline sm:block">{track.artistName}</a>
        </header>

        <div class="flex min-h-20 items-stretch">
            <div class="flex shrink-0 items-center border-r border-black/25 bg-black/10 px-2">
                <PlayPauseButton {track} />
            </div>

            <div
                class="relative flex min-w-0 flex-1 cursor-pointer items-center overflow-hidden bg-primary bg-[repeating-linear-gradient(90deg,rgba(0,0,0,0.16)_0,rgba(0,0,0,0.16)_1px,transparent_1px,transparent_48px)] px-2"
                bind:clientWidth={waveformWidth}
                onmousemove={handleMouseMove}
                onmouseenter={() => isCursorHovered = true}
                onmouseleave={() => isCursorHovered = false}
                onclick={playOrSeek}
                onkeydown={(event) => event.key === 'Enter' && playOrSeek()}
                role="button"
                tabindex="0"
                aria-label="Play or seek through {track.description}"
            >
                <div class="pointer-events-none absolute inset-y-0 w-px bg-black/45" style="left: {currentSongPercentage}%"></div>
                {#each visibleBars as bar, index}
                    <WaveformBar
                        height={bar}
                        {timePerBar}
                        {index}
                        {cursorTime}
                        {isCursorHovered}
                        playedPercentage={currentSongPercentage / 100}
                        trackDuration={track.duration}
                    />
                {/each}
                <div class="pointer-events-none absolute bottom-1 left-2 text-[10px] font-medium uppercase tracking-[0.16em] text-black/55 sm:hidden">{track.artistName} · {formatDuration(track.duration)}</div>
            </div>

            <aside class="grid w-20 shrink-0 grid-cols-2 content-center justify-items-center border-l border-black/25 bg-zinc-900/90 py-1 opacity-100 transition-opacity sm:opacity-55 sm:group-hover:opacity-100 sm:group-focus-within:opacity-100" aria-label="Project actions">
                <div class="relative flex flex-col items-center" title="Like">
                    <LikeButton trackId={track.id} bind:numLikes={likeCount} isLiked={isTrackLiked} showLabel={false} showCount={false} borderless={true} />
                    {#if likeCount > 0}<span class="absolute top-full text-[9px] text-zinc-400">{likeCount}</span>{/if}
                </div>

                <div class="relative" bind:this={shareArea}>
                    <button class="p-2 hover:cursor-pointer hover:bg-zinc-700 focus-visible:outline focus-visible:outline-1 focus-visible:outline-white" onclick={toggleSharePopup} aria-label="Share project">
                        <img src="send.png" alt="" class="h-4 w-4">
                    </button>
                    {#if showSharePopup}
                        <div class="absolute right-0 top-full z-50 mt-1 w-52 border border-zinc-700 bg-zinc-800 shadow-lg">
                            {#if copiedToClipboard}
                                <div class="px-4 py-2 text-sm text-white">Link copied</div>
                            {:else}
                                <button class="w-full cursor-pointer px-4 py-2 text-left text-sm text-white hover:bg-zinc-700" onclick={copyLink}>Copy link</button>
                            {/if}
                        </div>
                    {/if}
                </div>

                <AddToPageButton trackId={track.id} {urlBase} />
                <button class="p-2 hover:cursor-pointer hover:bg-zinc-700 focus-visible:outline focus-visible:outline-1 focus-visible:outline-white" onclick={navigateLayerr} aria-label="Use project">
                    <img src="vinyl.png" alt="" class="h-4 w-4">
                </button>
                {#if uses.length > 0}
                    <button class="col-span-2 border-t border-zinc-700 px-1 pt-1 text-[9px] font-semibold uppercase tracking-wide text-zinc-200 hover:text-white focus-visible:outline focus-visible:outline-1 focus-visible:outline-white" onclick={toggleStems} aria-expanded={isStemsExpanded}>
                        {isStemsExpanded ? 'Hide' : `Stems ${uses.length}`}
                    </button>
                {/if}
            </aside>
        </div>
    </div>

    {#if uses.length > 0 && isStemsExpanded}
        <section class="border-x border-b border-zinc-700 bg-zinc-950/80 py-2" aria-label="Source stems used by {track.description}">
            {#each uses as sourceTrack (sourceTrack.id)}
                <SourceStemLane track={sourceTrack} />
            {/each}
        </section>
    {/if}
</article>
