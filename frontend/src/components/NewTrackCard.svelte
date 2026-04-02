<script lang="ts">
    import { onMount, onDestroy, untrack } from 'svelte';
    import { goto } from '$app/navigation';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase } from '../stores/environment';
    import type { TrackInfo } from '../models/types';
    import { audio, currentTrack, currentTrackId, isPlaying, currentTime as globalCurrentTime } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import WaveformBar from './WaveformBar.svelte';
    import LikeButton from './LikeButton.svelte';
    
    type ColorName = 'red' | 'orange' | 'yellow' | 'green' | 'blue' | 'violet';

    const colorClasses: Record<ColorName, { text400: string; text500: string; hoverText500: string; border700: string; bg500: string; gradientBg: string; hoverShadow: string }> = {
        red:    { text400: 'text-red-400',    text500: 'text-red-500',    hoverText500: 'hover:text-red-500',    border700: 'border-red-700',    bg500: 'bg-red-500',    gradientBg: 'bg-gradient-to-r from-red-600 to-red-400',    hoverShadow: 'hover:shadow-red-500/30' },
        orange: { text400: 'text-orange-400', text500: 'text-orange-500', hoverText500: 'hover:text-orange-500', border700: 'border-orange-700', bg500: 'bg-orange-500', gradientBg: 'bg-gradient-to-r from-orange-600 to-orange-400', hoverShadow: 'hover:shadow-orange-500/30' },
        yellow: { text400: 'text-yellow-400', text500: 'text-yellow-500', hoverText500: 'hover:text-yellow-500', border700: 'border-yellow-700', bg500: 'bg-yellow-500', gradientBg: 'bg-gradient-to-r from-yellow-600 to-yellow-400', hoverShadow: 'hover:shadow-yellow-500/30' },
        green:  { text400: 'text-green-400',  text500: 'text-green-500',  hoverText500: 'hover:text-green-500',  border700: 'border-green-700',  bg500: 'bg-green-500',  gradientBg: 'bg-gradient-to-r from-green-600 to-green-400',  hoverShadow: 'hover:shadow-green-500/30' },
        blue:   { text400: 'text-blue-400',   text500: 'text-blue-500',   hoverText500: 'hover:text-blue-500',   border700: 'border-blue-700',   bg500: 'bg-blue-500',   gradientBg: 'bg-gradient-to-r from-blue-600 to-blue-400',   hoverShadow: 'hover:shadow-blue-500/30' },
        violet: { text400: 'text-violet-400', text500: 'text-violet-500', hoverText500: 'hover:text-violet-500', border700: 'border-violet-700', bg500: 'bg-violet-500', gradientBg: 'bg-gradient-to-r from-violet-600 to-violet-400', hoverShadow: 'hover:shadow-violet-500/30' },
    };

    // Inherits the track data from the page
    let { track }: { track: TrackInfo } = $props();

    const colors = $derived(colorClasses[(track.color as ColorName) ?? 'violet']);

    // State variables for the page
    let newAudioURL = $state('');

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
    let urlBase = $state('');

    // Waveform container width
    let waveformWidth = $state(0);
    let visibleBars = $state([0]);
    let timePerBar = $state(0);

    // Track data
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
                currentSongPercentage = currentTime / track.duration * 100
            }

            $audio.addEventListener('timeupdate', updateTime);

            return () => {
                $audio.removeEventListener('timeupdate', updateTime);
            }
        }
    })

    // When the component is loaded, sets up the environment and waveform
    onMount(async () => {
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
            visibleBars = selectedIndices.map(index => track.waveformData[index])
            timePerBar = track.duration / numBars
        }
    }

    function getSelectedIndices(numBars: number) {
            const selectedIndices = []
            const length = track.waveformData.length
            
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
        return currentTime / track.duration
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
        cursorTime = percentage * track.duration
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

<div class="w-2/3 max-w-[1200px] py-2 mb-4 border-b border-zinc-700">

    <div class="w-full h-12 mb-1 flex flex-row items-center">
            {#if track.artistPortraitUrl}
                <img src={track.artistPortraitUrl} alt={track.artistName} class="w-10 h-10 rounded-lg object-cover ml-2" />
            {/if}
            <a class="ml-2 px-1 rounded-md hover:bg-white text-lg transition-all duration-300 {colors.text500}" href={`/artist/${track.artistId}`}>{track.artistName}</a>
            <p class="ml-2 {colors.text400}">•</p>
            <a class="ml-2 px-1 rounded-md text-zinc-100 hover:bg-white text-lg transition-all duration-300 {colors.hoverText500}" href={`/track/${track.id}`}>{track.description}</a>
    </div>
    <!-- Waveform -->
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
                color={track.color}
            />
        {/each}
    </div>

    <!-- Track Information -->
    <div class="w-full h-12 flex flex-row items-center">
        <LikeButton trackId={track.id} numLikes={track.likes} isLiked={track.isLiked} color={track.color}></LikeButton>
        <button class="py-1 px-3 ml-4 rounded-md flex flex-row items-center justify-center transition-all duration-200 text-white font-semibold tracking-wider text-sm hover:scale-105 active:scale-95 hover:shadow-md {colors.gradientBg} {colors.hoverShadow}" onclick={navigateLayerr}>
            <p class="text-md">BUILD ON THIS</p>
        </button>
    </div>
</div>