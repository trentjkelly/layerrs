<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { goto } from '$app/navigation';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase } from '../stores/environment';
    import type { TrackInfo } from '../models/types';
    import { audio, currentTrack, currentTrackId, isPlaying, currentTime as globalCurrentTime } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import WaveformBar from './WaveformBar.svelte';

    let { track }: { track: TrackInfo } = $props();

    let stems = $state<TrackInfo[]>([]);

    let newAudioURL = $state('');
    let isExpanded = $state(false);
    let isTrackLiked = $state(false);
    let isHovered = $state(false);
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

            $audio.addEventListener('timeupdate', updateTime);

            return () => {
                $audio.removeEventListener('timeupdate', updateTime);
            };
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

    async function playPauseAudio() {
        if ($audio) {
            if (track.id === $currentTrackId) {
                if ($audio.paused) {
                    isPlaying.set(true);
                    await $audio.play();
                } else {
                    isPlaying.set(false);
                    $audio.pause();
                }
            } else {
                if (!$audio.paused) {
                    isPlaying.set(false);
                    $audio.pause();
                }

                if (mediaSource) {
                    if (sourceBuffer) {
                        try {
                            mediaSource.removeSourceBuffer(sourceBuffer);
                        } catch (error) {
                            logger.error(`Error removing source buffer: ${error}`);
                        }
                    }
                    mediaSource = null;
                    sourceBuffer = null;
                }

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
        }
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

    let playing = $derived($currentTrackId === track.id && $isPlaying);
</script>

<div class="w-full bg-olive-500 mb-2">
    <div class="flex flex-row h-16 items-center">
        <div class="mx-2 w-12 h-12">
            <img src={track.artistPortraitUrl} alt="Artist Profile" class="w-full h-full object-cover">
        </div>
        <div class="h-12 flex flex-col">
            <a href="/profile/{track.artistId}" class="text-md text-gray-300 hover:underline">{track.artistName}</a>
            <a href="/track/{track.id}" class="text-lg text-white hover:underline">{track.description}</a>
        </div>
        <div class="ml-auto mr-2 flex flex-row">
            <button class="mx-1 px-2 py-1 outline-1 hover:bg-olive-600 hover:cursor-pointer">Like</button>
            <button class="mx-1 px-2 py-1 outline-1 hover:bg-olive-600 hover:cursor-pointer flex flex-row items-center" onclick={navigateLayerr}>
                <img src="plus.png" alt="Layer" class="w-3 h-3">
                <p class="ml-1">Layer</p>
            </button>
            
            <div class="relative">
                <button 
                    id="share-btn"
                    class="mx-1 px-2 py-1 outline-1 hover:bg-olive-600 hover:cursor-pointer flex flex-row items-center"
                    onclick={toggleSharePopup}
                >
                    <img src="send.png" alt="Share" class="w-4 h-4">
                    <p class="ml-1">Share</p>
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
        </div>
    </div>

    <!-- Waveforms -->
    <div 
        class="relative h-16 w-full hover:cursor-pointer flex flex-row items-center py-1 px-1"
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

    <!-- Stems -->
    <div class="h-6 w-full bg-olive-600 flex flex-row items-center">
        <p class="mx-4">STEMS</p>

        {#each stems as stem}
            <div class="h-6 w-6 rounded-full bg-olive-600 mr-2"></div>
        {/each}
    </div>
</div>