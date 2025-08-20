<script lang="ts">
    import { onMount, onDestroy, untrack } from 'svelte';
    import { goto } from '$app/navigation';
    import { logger } from '../modules/lib/logger';
    import { getTrackData } from '../modules/requests/track-requests';
    import { getArtistName } from '../modules/requests/artist-requests';
    import { getUrlBase, handleEnvironment } from '../stores/environment';
    import { audio, currentTrack, currentTrackId, isPlaying, currentTime as globalCurrentTime } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import WaveformBar from './WaveformBar.svelte';
    import LikeButton from './LikeButton.svelte';
    import LayerrButton from './LayerrButton.svelte';
    
    // Inherits the trackId from the page
    let { trackId } = $props();

    // State variables for the page
    let newAudioURL = $state('');
    let trackName = $state('cannnot find track name');
    let artistId = $state(0);
    let artistName = $state('cannot find artist name');
    
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
    let numLikes = $state(0);
    let numLayerrs = $state(0);
    let urlBase = $state('');

    // Waveform container width
    let waveformWidth = $state(0);
    let waveformBars = $state([0]);
    let visibleBars = $state([0]);
    let timePerBar = $state(0);

    // Track data
    let trackDuration = $state(0);
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

    // When the component is loaded, gets the track data & cover art 
    onMount(async () => {
        await handleEnvironment()
        urlBase = getUrlBase()

        const trackData = await getTrackData(urlBase, trackId)
        if (trackData) {
            trackName = trackData.name
            artistId = parseInt(trackData.artistId)
            numLikes = trackData.likes
            numLayerrs = trackData.layerrs
            waveformBars = trackData.waveformData
            trackDuration = trackData.duration
        } else {
            console.log("Could not find track data")
        }

        const artistData = await getArtistName(urlBase, artistId)
        if (artistData) {
            artistName = artistData.name
        } else {
            console.log("Could not find artist data")
        }

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

    // Changes hover property when someone hovers the cover image
    function hoverTrackImage() {
        isHovered = true
    }

    // Changes hover property when someone unhovers the cover image
    function leaveHoverTrackImage() {
        isHovered = false
    }

    // Plays/pauses the audio
    async function playPauseAudio() {
        if ($audio) {
            // This Track is the current one (stored in session data)
            if (trackId === $currentTrackId) {
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
                newAudioURL = await getAudio(urlBase, trackId)
                currentTrack.set(newAudioURL)
                currentTrackId.set(trackId)

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
        goto(`/track/${trackId}`)
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
            if ($currentTrackId === trackId) {
                isPlaying.set(true)
                $audio.currentTime = cursorTime
                await $audio.play()
            } else {
                isPlaying.set(true)
                newAudioURL = await getAudio(urlBase, trackId)
                currentTrack.set(newAudioURL)
                currentTrackId.set(trackId)
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

<div class="w-3/4 max-w-[1200px] py-2 hover:cursor-pointer">

    <div class="w-full h-16 flex flex-row items-center">
        <div class="w-10 h-10 ml-4 rounded rounded-xl bg-gray-400"></div>
        <a class="ml-2 text-violet-400 hover:text-violet-500 hover:underline font-bold text-lg" href={`/artist/${artistId}`}>{artistName}</a>
        <p class="ml-2 text-violet-400">•</p>
        <a class="ml-2 text-white hover:text-violet-500 hover:underline text-md" href={`/track/${trackId}`}>{trackName}</a>
    </div>
    <!-- Waveform -->
    <div 
        class="relative h-36 w-full hover:cursor-pointer flex flex-row items-center rounded-2xl py-1 px-1 border     border-violet-700"
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
    <div class="w-full h-12 flex flex-row items-center px-4">
        <!-- Like Button -->
        <LikeButton trackId={trackId} numLikes={numLikes}></LikeButton>
        <LayerrButton trackId={trackId} numLayerrs={numLayerrs}></LayerrButton>
        <!-- Download Button -->
    </div>
</div>