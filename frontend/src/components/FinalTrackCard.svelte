<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { goto } from '$app/navigation';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase } from '../stores/environment';
    import type { TrackInfo, PageWithFollowerCount } from '../models/types';
    import { audio, currentTrack, currentTrackId, isPlaying, currentTime as globalCurrentTime } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import { trackPlayProgress } from '../modules/lib/play-tracking';
    import WaveformBar from './WaveformBar.svelte';
    import LikeButton from './LikeButton.svelte';
    import AddToPageButton from './AddToPageButton.svelte';
    import PlayPauseButton from './PlayPauseButton.svelte';
    import SourceStemLane from './SourceStemLane.svelte';

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

    let newAudioURL = $state('');
    let isExpanded = $state(false);
    let isTrackLiked = $state(track.isLiked);
    let likeCount = $state(track.likes);
    let showSharePopup = $state(false);
    let copiedToClipboard = $state(false);

    function toggleSharePopup(event: MouseEvent) {
        event.stopPropagation();
        showSharePopup = !showSharePopup;
        copiedToClipboard = false;
    }

    function closeSharePopup() {
        showSharePopup = false;
        copiedToClipboard = false;
    }

    function copyLink(event: MouseEvent) {
        event.stopPropagation();
        navigator.clipboard.writeText(`${window?.location?.origin}/track/${track.id}`);
        copiedToClipboard = true;
    }

    function handleClickOutside(event: MouseEvent) {
        const popup = document.getElementById('share-popup');
        const shareBtn = document.getElementById('share-btn');
        if (popup && !popup.contains(event.target as Node) && !shareBtn?.contains(event.target as Node)) {
            closeSharePopup();
        }
    }

    $effect(() => {
        if (showSharePopup) {
            document.addEventListener('click', handleClickOutside);
            return () => document.removeEventListener('click', handleClickOutside);
        }
    });

    let mediaSource = $state<MediaSource | null>(null);
    let sourceBuffer = $state<SourceBuffer | null>(null);
    let isLoading = $state(false);
    let urlBase = $state('');

    let waveformWidth = $state(0);
    let visibleBars = $state([0]);
    let timePerBar = $state(0);

    let currentTime = $state(0);
    let cursorTime = $state(0);
    let cursorPercentage = $state(0);
    let isCursorHovered = $state(false);
    let currentSongPercentage = $state(0);

    $effect(() => {
        changeWaveformWidth();
    });

    $effect(() => {
        if ($audio) {
            const updateTime = () => {
                if ($currentTrackId === track.id) {
                    currentTime = $audio.currentTime;
                    globalCurrentTime.set($audio.currentTime);
                    currentSongPercentage = currentTime / track.duration * 100;
                }
            };

            const handleEnded = () => {
                if ($currentTrackId === track.id) {
                    currentTime = 0;
                    globalCurrentTime.set(0);
                    currentSongPercentage = 0;
                }
            };

            $audio.addEventListener('timeupdate', updateTime);
            $audio.addEventListener('ended', handleEnded);

            return () => {
                $audio.removeEventListener('timeupdate', updateTime);
                $audio.removeEventListener('ended', handleEnded);
            };
        }
    });

    $effect(() => {
        if ($currentTrackId === track.id && urlBase) {
            return trackPlayProgress(track.id, urlBase);
        }
    });

    onMount(async () => {
        urlBase = getUrlBase();
        changeWaveformWidth();
    });

    onDestroy(() => {
        if (mediaSource) {
            if (sourceBuffer) {
                try {
                    mediaSource.removeSourceBuffer(sourceBuffer);
                } catch (error) {
                    logger.error(`Error removing the source buffer: ${error}`);
                }
            }
            mediaSource = null;
            sourceBuffer = null;
        }
        if (newAudioURL) {
            URL.revokeObjectURL(newAudioURL);
        }
    });

    function changeWaveformWidth() {
        const barWidth = 2;
        const barMargin = 2;
        const totalBarWidth = barWidth + barMargin;
        const numBars = Math.floor(waveformWidth / totalBarWidth);

        if (numBars <= 0 || track.waveformData.length === 0) {
            visibleBars = [];
            timePerBar = 0;
            return;
        }

        const selectedIndices = getSelectedIndices(numBars);
        visibleBars = selectedIndices.map(index => track.waveformData[index]);
        timePerBar = track.duration / numBars;
    }

    function getSelectedIndices(numBars: number) {
        const selectedIndices: number[] = [];
        const length = track.waveformData.length;

        for (let i = 0; i < numBars; i++) {
            const index = Math.min(length - 1, Math.floor(length * i / numBars));
            selectedIndices.push(index);
        }
        return selectedIndices;
    }

    function onLikeAndDownload() {
        isExpanded = true;
    }

    function offLikeAndDownload() {
        isExpanded = false;
    }

    function navigateLayerr() {
        goto(`/layerrs/${track.id}`, { state: { track: $state.snapshot(track) } });
    }

    function getSongPercentage() {
        return currentTime / track.duration;
    }

    function handleMouseMove(event: MouseEvent) {
        const cursorX = event.clientX;
        const left = (event.currentTarget as HTMLElement).getBoundingClientRect().left;
        let percentage = (cursorX - left) / waveformWidth;
        if (percentage < 0) {
            percentage = 0;
        } else if (percentage > 1) {
            percentage = 1;
        }
        cursorPercentage = percentage * 100;
        cursorTime = percentage * track.duration;
    }

    function handleMouseLeave() {
        isCursorHovered = false;
    }

    function handleMouseEnter() {
        isCursorHovered = true;
    }

    async function handleClick(event: MouseEvent) {
        if ($audio) {
            if ($currentTrackId === track.id) {
                isPlaying.set(true);
                $audio.currentTime = cursorTime;
                await $audio.play();
            } else {
                isPlaying.set(true);
                newAudioURL = await getAudio(urlBase, String(track.id));
                currentTrack.set(newAudioURL);
                currentTrackId.set(track.id);
                $audio.src = newAudioURL;
                try {
                    await $audio.play();
                } catch (error) {
                    logger.error(`Failed to play audio: ${error}`);
                }
            }
        } else {
            console.error("Failed to load audio element");
        }
    }

    function formatPlays(n: number): string {
        if (n >= 1_000_000) return (n / 1_000_000).toFixed(1).replace(/\.0$/, '') + 'm';
        if (n >= 1_000) return (n / 1_000).toFixed(1).replace(/\.0$/, '') + 'k';
        return String(n);
    }

    function formatDuration(seconds: number): string {
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    }

</script>

<article class="w-full overflow-hidden border-b border-zinc-700 bg-zinc-900">
    <div class="flex min-h-28 w-full flex-col sm:min-h-24 sm:flex-row">
        <div class="h-1 w-full shrink-0 bg-orange-500 sm:h-auto sm:w-1"></div>

        <header class="flex min-w-0 items-center gap-3 px-3 py-3 sm:w-64 sm:shrink-0 sm:border-r sm:border-zinc-700">
            <span class="w-6 shrink-0 text-center text-xs text-orange-400">
                {String(trackNumber || track.id).padStart(2, '0')}
            </span>
            <img src={track.artistPortraitUrl} alt="{track.artistName} profile" class="h-10 w-10 shrink-0 object-cover">
            <div class="min-w-0 flex-1">
                <a href="/track/{track.id}" class="block truncate text-sm font-medium text-white hover:underline">{track.description}</a>
                <a href="/profile/{track.artistId}" class="block truncate text-xs text-zinc-400 hover:text-white hover:underline">{track.artistName}</a>
                <div class="mt-1 flex items-center gap-2 text-[10px] uppercase tracking-wide text-zinc-500">
                    <span>{formatDuration(track.duration)}</span>
                    {#if track.plays > 0}
                        <span>{formatPlays(track.plays)} plays</span>
                    {/if}
                </div>
            </div>
        </header>

        <div class="flex min-w-0 flex-1 items-stretch border-t border-zinc-700 sm:border-t-0">
            <div class="flex shrink-0 items-center border-r border-zinc-700 px-2">
                <PlayPauseButton {track} />
            </div>
            <div
                class="relative flex h-20 min-w-0 flex-1 cursor-pointer items-center overflow-hidden bg-orange-500 bg-[repeating-linear-gradient(90deg,rgba(0,0,0,0.16)_0,rgba(0,0,0,0.16)_1px,transparent_1px,transparent_48px)] px-2"
                bind:clientWidth={waveformWidth}
                onmousemove={handleMouseMove}
                onmouseleave={handleMouseLeave}
                onmouseenter={handleMouseEnter}
                onclick={handleClick}
                onkeydown={(e) => e.key === 'Enter' && handleClick(e as any)}
                role="button"
                tabindex="0"
                aria-label="Seek through {track.description}"
            >
                <div class="pointer-events-none absolute inset-y-0 left-0 w-px bg-black/40" style="left: {currentSongPercentage}%"></div>
                {#each visibleBars as bar, index}
                    <WaveformBar
                        height={bar}
                        timePerBar={timePerBar}
                        index={index}
                        cursorTime={cursorTime}
                        isCursorHovered={isCursorHovered}
                        playedPercentage={currentSongPercentage / 100}
                        trackDuration={track.duration}
                    />
                {/each}
            </div>
        </div>

        <div class="flex shrink-0 items-center justify-end gap-1 border-t border-zinc-700 px-2 py-2 sm:w-44 sm:border-t-0 sm:border-l sm:border-zinc-700">
            <div class="relative flex flex-col items-center justify-center">
                <LikeButton trackId={track.id} bind:numLikes={likeCount} isLiked={isTrackLiked} showLabel={false} showCount={false} borderless={true} />
                {#if likeCount > 0}
                    <span class="absolute top-full text-[10px] text-zinc-400">{likeCount}</span>
                {/if}
            </div>

            <div class="relative flex flex-col items-center justify-center">
                <button id="share-btn" class="p-2 hover:cursor-pointer hover:bg-zinc-700" onclick={toggleSharePopup} aria-label="Share">
                    <img src="send.png" alt="" class="h-4 w-4">
                </button>
                {#if showSharePopup}
                    <div id="share-popup" class="absolute right-0 top-full z-50 mt-1 w-52 border border-zinc-700 bg-zinc-800">
                        {#if copiedToClipboard}
                            <div class="flex items-center px-4 py-2 text-white">
                                <svg class="mr-3 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
                                <span class="text-sm">Added to Clipboard</span>
                            </div>
                        {:else}
                            <button class="flex w-full cursor-pointer items-center px-4 py-2 text-left text-sm text-white transition-colors hover:bg-zinc-700" onclick={copyLink}>
                                <svg class="mr-3 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path></svg>
                                Copy Link
                            </button>
                        {/if}
                    </div>
                {/if}
            </div>

            <AddToPageButton trackId={track.id} {urlBase} />
            <button class="flex items-center gap-1 border border-zinc-600 px-2 py-1 text-xs hover:cursor-pointer hover:border-zinc-400 hover:bg-zinc-700" onclick={navigateLayerr}>
                <img src="vinyl.png" alt="" class="h-4 w-4">
                <span>USE</span>
            </button>
        </div>
    </div>

    {#if uses.length > 0}
        <section class="border-t border-zinc-800 bg-zinc-950/40 py-2" aria-label="Source stems used by {track.description}">
            {#each uses as sourceTrack (sourceTrack.id)}
                <SourceStemLane track={sourceTrack} />
            {/each}
        </section>
    {/if}
</article>
