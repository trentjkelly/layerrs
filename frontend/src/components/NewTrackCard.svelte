<script lang="ts">
    import { onMount, onDestroy, untrack } from 'svelte';
    import { goto } from '$app/navigation';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase, handleEnvironment } from '../stores/environment';
    import type { TrackInfo } from '../models/types';
    import { audio, currentTrack, currentTrackId, isPlaying, currentTime as globalCurrentTime } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import WaveformBar from './WaveformBar.svelte';
    import LikeButton from './LikeButton.svelte';
    
    type ColorName = 'red' | 'orange' | 'yellow' | 'green' | 'blue' | 'violet';

    const colorClasses: Record<ColorName, { text400: string; text500: string; hoverText500: string; border700: string }> = {
        red:    { text400: 'text-red-400',    text500: 'text-red-500',    hoverText500: 'hover:text-red-500',    border700: 'border-red-700' },
        orange: { text400: 'text-orange-400', text500: 'text-orange-500', hoverText500: 'hover:text-orange-500', border700: 'border-orange-700' },
        yellow: { text400: 'text-yellow-400', text500: 'text-yellow-500', hoverText500: 'hover:text-yellow-500', border700: 'border-yellow-700' },
        green:  { text400: 'text-green-400',  text500: 'text-green-500',  hoverText500: 'hover:text-green-500',  border700: 'border-green-700' },
        blue:   { text400: 'text-blue-400',   text500: 'text-blue-500',   hoverText500: 'hover:text-blue-500',   border700: 'border-blue-700' },
        violet: { text400: 'text-violet-400', text500: 'text-violet-500', hoverText500: 'hover:text-violet-500', border700: 'border-violet-700' },
    };

    // Inherits the track data from the page
    let { track }: { track: TrackInfo } = $props();

    const colors = $derived(colorClasses[(track.color as ColorName) ?? 'violet']);

    // State variables for the page
    let newAudioURL = $state('');
    let trackDescription = $state(track.description);
    let artistId = $state(track.artistId);
    let artistName = $state(track.artistName);
    
    let parentTrackName = $state('');
    let parentTrackId = $state(0);
    let parentTrackArist = $state('Yer');

    let isExpanded = $state(false);
    let isTrackLiked = $state(false);
    let isHovered = $state(false);
    let audioElement = $state();
    let mediaSource = $state<MediaSource | null>(null);
    let sourceBuffer = $state<SourceBuffer | null>(null);
    let isLoading = $state(false);
    let currentOffset = $state(0);
    let numLikes = $state(track.likes);
    let urlBase = $state('');

    // Waveform container width
    let waveformWidth = $state(0);
    let waveformBars = $state(track.waveformData);
    let visibleBars = $state([0]);
    let timePerBar = $state(0);

    // Track data
    let trackDuration = $state(track.duration);
    let currentTime = $state(0);
    let cursorTime = $state(0);
    let cursorPercentage = $state(0);
    let isCursorHovered = $state(false);

    let currentSongPercentage = $state(0);
    
    // Reactive statement to generate waveform bars based on width
    $effect(() => {
        changeWaveformWidth()
    });

    $effect(() => {
        if ($audio) {
            const updateTime = () => {
                currentTime = $audio.currentTime
                globalCurrentTime.set($audio.currentTime)
                currentSongPercentage = currentTime / trackDuration * 100
            }

            $audio.addEventListener('timeupdate', updateTime);

            return () => {
                $audio.removeEventListener('timeupdate', updateTime);
            }
        }
    })

    // When the component is loaded, sets up the environment and waveform
    onMount(async () => {
        await handleEnvironment()
        urlBase = getUrlBase()
        changeWaveformWidth()
    })

    // Added cleanup when component is destroyed
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
            sourceBuffer = null
        }
        if(newAudioURL) {
            URL.revokeObjectURL(newAudioURL);
        }
    });

    function changeWaveformWidth() {
        if (waveformWidth > 0) {
            const barWidth = 2;
            const barMargin = 2;
            const totalBarWidth = barWidth + barMargin;
            const numBars = Math.floor(waveformWidth / totalBarWidth);
            const selectedIndices = getSelectedIndices(numBars)
            visibleBars = selectedIndices.map(index => waveformBars[index])
            timePerBar = trackDuration / numBars
        }
    }

    function getSelectedIndices(numBars: number) {
            const selectedIndices = []
            const length = waveformBars.length
            
            for (let i = 0; i < numBars; i++) {
                let index = Math.round(length * i / numBars)
                selectedIndices.push(index)
            }
            return selectedIndices
    }

    // Plays/pauses the audio
    async function playPauseAudio() {
        if ($audio) {
            // This Track is the current one (stored in session data)
            if (track.id === $currentTrackId) {
                if ($audio.paused) {
                    isPlaying.set(true)
                    await $audio.play()
                } else {
                    isPlaying.set(false)
                    $audio.pause()
                }
            } 
            // This track is different than the current one (stored in session data)
            else {
                // Pause current audio
                if (!$audio.paused) {
                    isPlaying.set(false)
                    $audio.pause()
                }

                // Clean up old MediaSource if it exists
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

                // Play new audio
                isPlaying.set(true)
                newAudioURL = await getAudio(urlBase, String(track.id))
                currentTrack.set(newAudioURL)
                currentTrackId.set(track.id)

                $audio.src = newAudioURL
                try {
                    await $audio.play()
                } catch (error) {
                    logger.error(`Failed to play audio: ${error}`);
                }
            }
        }
    }

    // Shows the like and download buttons when hovered over
    function onLikeAndDownload() {
        isExpanded = true
    }

    // Hides the like and download buttons when hover is left
    function offLikeAndDownload() {
        isExpanded = false
    }

    function navigateTrackPage() {
        goto(`/track/${track.id}`)
    }

    function navigateLayerr() {
        console.log("I'm here")
        goto(`/layerrs/${track.id}`, { state: { track: $state.snapshot(track) } })
    }

    function getSongPercentage() {
        return currentTime / trackDuration
    }

    function handleMouseMove(event: MouseEvent) {
        const cursorX = event.clientX
        const left = (event.currentTarget as HTMLElement).getBoundingClientRect().left
        let percentage = (cursorX - left) / waveformWidth
        if (percentage < 0) {
            percentage = 0
        } else if (percentage > 1) {
            percentage = 1
        }
        cursorPercentage = percentage * 100
        cursorTime = percentage * trackDuration
    }

    function handleMouseLeave(event: MouseEvent) {
        isCursorHovered = false
    }

    function handleMouseEnter(event: MouseEvent) {
        isCursorHovered = true
    }

    async function handleClick(event: MouseEvent) {
        if($audio) {
            if ($currentTrackId === track.id) {
                isPlaying.set(true)
                $audio.currentTime = cursorTime
                await $audio.play()
            } else {
                isPlaying.set(true)
                newAudioURL = await getAudio(urlBase, String(track.id))
                currentTrack.set(newAudioURL)
                currentTrackId.set(track.id)
                $audio.src = newAudioURL
                try {
                    await $audio.play()
                } catch (error) {
                    logger.error(`Failed to play audio: ${error}`);
                }
            }
        } else {
            console.error("Failed to load audio element")
        }
    }

</script>

<div class="w-3/4 max-w-[1200px] py-2">

    <div class="w-full h-8 mb-1 flex flex-row items-center">
            {#if track.artistPortraitUrl}
                <img src={track.artistPortraitUrl} alt={artistName} class="w-7 h-7 rounded-lg object-cover ml-2" />
            {/if}
            <a class="ml-2 px-1 hover:bg-white text-lg transition-all duration-300 {colors.text500}" href={`/artist/${artistId}`}>{artistName}</a>
            <p class="ml-2 {colors.text400}">•</p>
            <a class="ml-2 px-1 text-zinc-100 hover:bg-white text-lg transition-all duration-300 {colors.hoverText500}" href={`/track/${track.id}`}>{trackDescription}</a>
    </div>
    <!-- Waveform -->
    <div 
        class="relative h-24 w-full hover:cursor-pointer flex flex-row items-center rounded-2xl py-1 px-1 border-2 {colors.border700}"
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
                color={track.color}
            />
        {/each}
    </div>

    <!-- Track Information -->
    <div class="w-full h-12 flex flex-row items-center">
        <LikeButton trackId={track.id} numLikes={numLikes} isLiked={track.isLiked} color={track.color}></LikeButton>
        <button class="py-1 px-2 ml-4 flex flex-row items-center justify-center hover:bg-white transition-all duration-300 text-white {colors.hoverText500}" onclick={navigateLayerr}>
            <p class="text-md">BUILD ON THIS</p>
        </button>
    </div>
</div>