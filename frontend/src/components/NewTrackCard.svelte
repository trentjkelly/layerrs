<script lang="ts">
    import { onMount, onDestroy, untrack } from 'svelte';
    import { goto } from '$app/navigation';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase } from '../stores/environment';
    import type { TrackInfo } from '../models/types';
    import { audio, currentTrack, currentTrackId, isPlaying, currentTime as globalCurrentTime } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import { trackPlayProgress } from '../modules/lib/play-tracking';
    import WaveformBar from './WaveformBar.svelte';
    import LikeButton from './LikeButton.svelte';

    // Inherits the track data from the page
    let { track }: { track: TrackInfo } = $props();

    // State variables for the page
    let newAudioURL = $state('');

    let parentTrackName = $state('');
    let parentTrackId = $state(0);
    let parentTrackArist = $state('Yer');

    let isExpanded = $state(false);
    let isTrackLiked = $state(false);
    let numLikes = $state(track.likes);
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

    $effect(() => {
        if ($currentTrackId === track.id && urlBase) {
            return trackPlayProgress(track.id, urlBase);
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
            <a class="ml-2 px-1 rounded-md hover:bg-white text-lg transition-all duration-300" href={`/artist/${track.artistId}`}>{track.artistName}</a>
            <p class="ml-2">•</p>
            <a class="ml-2 px-1 rounded-md text-zinc-100 hover:bg-white text-lg transition-all duration-300" href={`/track/${track.id}`}>{track.description}</a>
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
            />
        {/each}
    </div>

    <!-- Track Information -->
    <div class="w-full h-12 flex flex-row items-center">
        <LikeButton trackId={track.id} bind:numLikes={numLikes} isLiked={track.isLiked}></LikeButton>
        <button class="py-1 px-3 ml-4 rounded-md flex flex-row items-center justify-center transition-all duration-200 text-white font-semibold tracking-wider text-sm hover:scale-105 active:scale-95 hover:shadow-md" onclick={navigateLayerr}>
            <p class="text-md">BUILD ON THIS</p>
        </button>
    </div>
</div>