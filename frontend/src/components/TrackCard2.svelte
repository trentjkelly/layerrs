<script>
    import { onMount, onDestroy } from 'svelte';
    import { audio, currentTrack, currentTrackId, isPlaying } from '../stores/player';
    import { isLoggedIn, jwt } from '../stores/auth';    
    import { goto } from '$app/navigation';
    import LikeButton from './LikeButton.svelte';
    import LayerrButton from './LayerrButton.svelte';
    import ParentTrackButton from './ParentTrackButton.svelte';
    import { urlBase } from '../stores/environment';
    import { logger } from '../lib/logger/logger';
    // Inherits the trackId from the page
    let { trackId } = $props();

    // State variables for the page
    let coverURL = $state('');
    let newAudioURL = $state('');
    let trackName = $state('yessir');
    let artistId = $state(0);
    let artistName = $state('');
    
    let parentTrackName = $state('');
    let parentTrackId = $state(0);
    let parentTrackArist = $state('Yer')

    let isExpanded = $state(false);
    let isTrackLiked = $state(false);
    let isHovered = $state(false);
    let audioElement = $state();
    let mediaSource = $state();
    let sourceBuffer = $state();
    let isLoading = $state(false);
    let currentOffset = $state(0);
    let numLikes = $state(0);
    let numLayerrs = $state(0);
    const CHUNK_SIZE = 1024 * 1024; // 1 MB

    // When the component is loaded -- gets the track data & cover art 
    onMount(async () => {
        await getTrackData()
        await getArtistName()
        await getCover()
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

    // Requests the metadata for the track
    async function getTrackData() {
        try {
            if ($urlBase) {
                const response = await fetch(`${$urlBase}/api/track/${trackId}/data`, { method: "GET"});
                if (!response.ok) {
                    throw new Error("Failed to get track data");
                }
                const trackData = await response.json();
                trackName = trackData.name
                artistId = trackData.artistId
                numLikes = trackData.likes
                numLayerrs = trackData.layerrs
            }
        } catch (error) {
            logger.error(`Error catching track data: ${error}`);
        }
    }

    // Gets the name of the artist
    async function getArtistName() {
        try {
            // const backendURL = import.meta.env.VITE_BACKEND_URL;
            const response = await fetch(`${$urlBase}/api/artist/${artistId}`, {
                method: "GET"
            })
            if(!response.ok) {
                throw new Error("Failed to get artist data");
            }

            const artistData = await response.json();
            artistName = artistData.name
        } catch (error) {
            logger.error(`Could not retrieve artist data: ${error}`);
        }
    }

    // Requests the cover art for the track
    async function getCover() {
        try {
            // const backendURL = import.meta.env.VITE_BACKEND_URL;
            const response = await fetch(`${$urlBase}/api/track/${trackId}/cover`, { method: "GET"});
            if (!response.ok) {
                throw new Error("Failed to catch cover art");
            }
            const blob = await response.blob();
            coverURL = URL.createObjectURL(blob);

        } catch (error) {
            logger.error(`Error catching cover art: ${error}`);
        }
    }

    // Changes hover property when someone hovers the cover image
    function hoverTrackImage() {
        isHovered = true
    }

    // Changes hover property when someone unhovers the cover image
    function leaveHoverTrackImage() {
        isHovered = false
    }

    // Requests the audio for the track
    async function getAudio() {
        try {
            mediaSource = new MediaSource();
            audioElement = new Audio();
            currentOffset = 0;

            const sourceURL = URL.createObjectURL(mediaSource);
            audioElement.src = sourceURL;

            mediaSource.addEventListener('sourceopen', async () => {
                try {
                    sourceBuffer = mediaSource.addSourceBuffer('audio/mpeg');
                    
                    if (sourceBuffer) {
                        sourceBuffer.addEventListener('updateend', () => {
                            if (!isLoading) {
                                loadNextChunk();
                            }
                        });
                    }
                    await loadNextChunk();
                } catch (error) {
                    logger.error(`Error setting up media source: ${error}`);
                }
            });

            audio.set(audioElement);
            newAudioURL = sourceURL;

        } catch (error) {
            logger.error(`Error setting up audio stream: ${error}`);
        }
    }

    async function loadNextChunk() {
        if (isLoading || !mediaSource || mediaSource.readyState !== 'open') {
            return
        }

        try {
            isLoading = true;
            // const backendURL = import.meta.env.VITE_BACKEND_URL;
            const response = await fetch(`${$urlBase}/api/track/${trackId}/audio`,
            {
                headers: {
                    'Range': `bytes=${currentOffset}-${currentOffset + CHUNK_SIZE - 1}`
                }
            });

            if(!response.ok) {
                throw new Error('Failed to fetch chunk');
            }

            const data = await response.arrayBuffer();
            if(data.byteLength === 0) {
                mediaSource.endOfStream();
                return;
            }

            if (!sourceBuffer.updating) {
                sourceBuffer.appendBuffer(data);
                currentOffset += data.byteLength;
            }

        } catch (error) {
            logger.error(`Error loading chunk: ${error}`);
        } finally {
            isLoading = false;
        }
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
                await getAudio()
                currentTrack.set(newAudioURL)
                currentTrackId.set(trackId)

                // if ($currentTrack) {
                //     $audio.src = $currentTrack
                // }
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

</script>

<div class="w-64 h-auto flex flex-col justify-center mb-4 mx-1 bg-slate-900 hover:bg-slate-800" onmouseenter={onLikeAndDownload} onfocus={onLikeAndDownload} onmouseleave={offLikeAndDownload} role="button" tabindex="0">
    <!-- Picture section -->
    <div class="h-64 w-64 flex flex-row items-center justify-center">
        <div 
            id={trackId}
            onmouseover={hoverTrackImage} 
            onfocus={hoverTrackImage}
            onmouseleave={leaveHoverTrackImage}
            onclick={playPauseAudio} 
            onkeydown={(e) => {if (e.key === "Enter" || e.key === " ") playPauseAudio}} 
            role="button" 
            tabindex="0" 
            class="h-64 w-64 bg-slate-700 flex flex-row items-center justify-center"
        >
            {#if coverURL}
                <img class="h-64 w-64 absolute" src={coverURL} alt="cover art">
            {/if}
            
            <!-- Absolutely disgusting conditional for whether to show play or pause button -->
            {#if isHovered}
                {#if $isPlaying}
                    {#if (trackId === $currentTrackId)}
                        <img class="h-20 w-20 absolute" src="/pause.png" alt="Pause button" />
                    {:else}
                        <img class="h-20 w-20 absolute" src="/play.png" alt="Play button" />
                    {/if}
                {:else}
                    <img class="h-20 w-20 absolute" src="/play.png" alt="Play button" />
                {/if}
            {:else}
                {#if ($isPlaying && (trackId === $currentTrackId))}
                    <img class="h-20 w-20 absolute" src="/pause.png" alt="Pause button" />
                {/if}
            {/if}
            
        </div>
    </div>

    <!-- Section below the picture -->  
    <div class="w-full h-28 px-4">
        <div class="flex flex-row w-full mt-2">
            <div class={`flex flex-col ${isExpanded ? 'w-40' : 'w-full'}`}>
                <a class="font-bold text-lg hover:underline truncate" href="/track/{trackId}">{trackName}</a>
                <a class="pb-2 text-gray-400 text-md hover:underline truncate" href="/artist/{artistId}">@{artistName}</a>
            </div>
            {#if isExpanded}
                <div class="w-28 flex flex-row items-center justify-between">
                    <LikeButton trackId={trackId} numLikes={numLikes}/>
                    <LayerrButton trackId={trackId} numLayerrs={numLayerrs} />
                </div>
            {/if}
        </div>
        <div class="pb-4">
            {#if parentTrackName}
                <ParentTrackButton parentTrackId={parentTrackId} parentTrackArist={parentTrackArist} parentTrackName={parentTrackName}/>
            {:else}
                <p class="text-blue-500">Original Track</p>
            {/if}
        </div>
    </div>
</div>
