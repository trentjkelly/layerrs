<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { goto } from '$app/navigation';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase } from '../stores/environment';
    import type { TrackInfo, PageWithFollowerCount } from '../models/types';
    import { audio, currentTrack, currentTrackId, isPlaying, currentTime as globalCurrentTime } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import { trackPlayProgress } from '../modules/lib/play-tracking';
    import { formatFollowers } from '../modules/lib/format';
    import WaveformBar from './WaveformBar.svelte';
    import LikeButton from './LikeButton.svelte';
    import AddToPageButton from './AddToPageButton.svelte';
    import PlayPauseButton from './PlayPauseButton.svelte';

    let {
        track,
        pages = [],
        pageCount = 0,
        uses = []
    }: {
        track: TrackInfo;
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
        if (waveformWidth > 0) {
            const barWidth = 2;
            const barMargin = 2;
            const totalBarWidth = barWidth + barMargin;
            const numBars = Math.floor(waveformWidth / totalBarWidth);
            const selectedIndices = getSelectedIndices(numBars);
            visibleBars = selectedIndices.map(index => track.waveformData[index]);
            timePerBar = track.duration / numBars;
        }
    }

    function getSelectedIndices(numBars: number) {
        const selectedIndices = [];
        const length = track.waveformData.length;

        for (let i = 0; i < numBars; i++) {
            let index = Math.round(length * i / numBars);
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

    function navigateTrackPage() {
        goto(`/track/${track.id}`);
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

<div class="w-full bg-olive-500 mb-2 shadow-lg">
    <div class="flex flex-row h-16 items-start pt-2">
        <div class="mx-2 w-12 h-12 relative flex flex-row items-center justify-center">
            <img src={track.artistPortraitUrl} alt="Artist Profile" class="w-full h-full object-cover">
        </div>
        <div class="h-12 flex flex-col">
            <div class="flex flex-row items-center gap-2">
                <a href="/track/{track.id}" class="text-lg text-white hover:underline">{track.description}</a>
                {#if track.plays > 0}
                    <span class="text-sm text-gray-300 flex flex-row items-center gap-1 ml-2">
                        <img src="play.png" alt="Plays" class="w-3 h-3">
                        {formatPlays(track.plays)}
                    </span>
                {/if}
            </div>
            <a href="/profile/{track.artistId}" class="text-sm text-gray-300 hover:underline">{track.artistName}</a>
        </div>
        <div class="ml-auto mr-2 flex flex-row items-center">
            <div class="relative flex flex-col items-center justify-center">
                <LikeButton
                    trackId={track.id}
                    bind:numLikes={likeCount}
                    isLiked={isTrackLiked}
                    showLabel={false}
                    showCount={false}
                    borderless={true}
                />
                {#if likeCount > 0}
                    <span class="absolute top-full text-xs text-gray-300">{likeCount}</span>
                {/if}
            </div>
            
            <div class="relative flex flex-col items-center justify-center">
                <button 
                    id="share-btn"
                    class="mr-1 p-2 hover:bg-olive-600 hover:cursor-pointer flex flex-row items-center"
                    onclick={toggleSharePopup}
                >
                    <img src="send.png" alt="Share" class="w-5 h-5">
                </button>
                
                {#if showSharePopup}
                    <div id="share-popup" class="absolute top-full right-0 mt-1 w-52 bg-olive-500 shadow-xl border border-white z-50">
                        {#if copiedToClipboard}
                            <div class="px-4 py-2 flex items-center text-white">
                                <svg class="w-4 h-4 mr-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                                </svg>
                                <span class="text-md">Added to Clipboard</span>
                            </div>
                        {:else}
                            <button 
                                class="w-full px-4 py-2 text-left text-white hover:bg-olive-600 flex items-center transition-colors duration-150 cursor-pointer"
                                onclick={copyLink}
                            >
                                <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path>
                                </svg>
                                Copy Link
                            </button>
                        {/if}
                    </div>
                {/if}
            </div>

            <AddToPageButton trackId={track.id} {urlBase} />

            <button class="mr-1 px-2 py-1 outline-1 hover:bg-olive-600 hover:cursor-pointer flex flex-row items-center" onclick={navigateLayerr}>
                <img src="vinyl.png" alt="Use" class="w-5 h-5">
                <p class="ml-1">Use</p>
            </button>
        </div>
    </div>

    <!-- Appears on -->
    {#if pageCount > 0}
        <div class="px-2 pb-2 flex flex-row items-center gap-2 text-sm flex-wrap">
            <span class="text-gray-300">Appears on:</span>
            {#each pages as page}
                <a
                    href="/pages/{page.id}"
                    class="bg-olive-600 hover:bg-olive-700 px-2 py-0.5 transition-colors"
                >
                    {page.name}
                    {#if (page.followerCount ?? 0) > 1}
                        <span class="text-gray-300">· {formatFollowers(page.followerCount ?? 0)}</span>
                    {/if}
                </a>
            {/each}
            {#if pageCount > pages.length}
                <button
                    type="button"
                    onclick={navigateTrackPage}
                    class="text-gray-300 hover:text-white hover:underline cursor-pointer transition-colors"
                >
                    +{pageCount - pages.length} more
                </button>
            {/if}
        </div>
    {/if}

    <!-- Track Duration -->
    <div class="flex flex-row justify-end px-1">
        <span class="text-sm text-gray-300 whitespace-nowrap">{formatDuration(track.duration)}</span>
    </div>

    <!-- Play/Pause + Waveforms -->
    <div class="flex flex-row items-center px-1 pt-0">
        <div class="mx-2">
            <PlayPauseButton {track} />
        </div>
        <div
            class="relative h-16 flex-1 hover:cursor-pointer flex flex-row items-center pt-0 pb-1"
            bind:clientWidth={waveformWidth}
            onmousemove={handleMouseMove}
            onmouseleave={handleMouseLeave}
            onmouseenter={handleMouseEnter}
            onclick={handleClick}
            onkeydown={(e) => e.key === 'Enter' && handleClick(e as any)}
            role="button"
            tabindex="0"
        >

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

    <!-- Uses -->
    {#if uses.length > 0}
        <div class="min-h-6 w-full bg-olive-600 flex flex-row items-center flex-wrap px-2 py-1">
            <p class="mx-2 text-sm">USES:</p>

            {#each uses as use}
                <a
                    href="/track/{use.id}"
                    class="flex flex-row items-center mr-2 mb-1 hover:underline"
                >
                    {#if use.artistPortraitUrl}
                        <img
                            src={use.artistPortraitUrl}
                            alt={use.artistName}
                            class="h-5 w-5 rounded-full object-cover mr-1"
                        />
                    {:else}
                        <div class="h-5 w-5 rounded-full bg-olive-700 mr-1"></div>
                    {/if}
                    <span class="text-sm">{use.artistName} - {use.description}</span>
                </a>
            {/each}
        </div>
    {/if}
</div>