<script>
    import { onMount, onDestroy } from 'svelte';
    import { audio, currentTrack, currentTrackId, isPlaying } from '../stores/player';
        
    // Inherits the trackId from the page
    let { trackId } = $props();

    // State variables for the page
    let coverURL = $state('');
    let newAudioURL = $state('');
    let trackName = $state('');
    let previousTrackName = $state('');
    let isExpanded = $state(false);
    let isTrackLiked = $state(false);
    let isHovered = $state(false);

    // New variables
    let audioElement = $state();
    let mediaSource = $state();
    let sourceBuffer = $state();
    let isLoading = $state(false);
    let currentOffset = $state(0);
    const CHUNK_SIZE = 1024 * 1024; // 1 MB

    // When the component is loaded -- gets the track data & cover art 
    onMount(async () => {
        await getTrackData()
        await getCover()
    })

    // Added cleanup when component is destroyed
    onDestroy(() => {
        if (mediaSource) {
            if (sourceBuffer) {
                try {
                    mediaSource.removeSourceBuffer(sourceBuffer);
                } catch (error) {
                    console.warn("Error removing the source buffer", error);
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
            const response = await fetch(`http://localhost:8080/api/track/${trackId}/data`, { method: "GET"});
            if (!response.ok) {
                throw new Error("Failed to get track data");
            }
            const trackData = await response.json();
            trackName = trackData.name
            previousTrackName = "Yessirski"
            // artistName = trackData.
        } catch (error) {
            console.error("Error catching track data", error)
        }
    }

    // Requests the cover art for the track
    async function getCover() {
        try {
            const response = await fetch(`http://localhost:8080/api/track/${trackId}/cover`, { method: "GET"});
            if (!response.ok) {
                throw new Error("Failed to catch cover art");
            }
            const blob = await response.blob();
            coverURL = URL.createObjectURL(blob);

        } catch (error) {
            console.error("Error catching cover art", error)
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
                    console.error("Error setting up media source", error);
                }
            });

            audio.set(audioElement);
            newAudioURL = sourceURL;

        } catch (error) {
            console.error("Error setting up audio stream", error);
        }
    }

    async function loadNextChunk() {
        if (isLoading || !mediaSource || mediaSource.readyState !== 'open') {
            return
        }

        try {
            isLoading = true;
            const response = await fetch(`http://localhost:8080/api/track/${trackId}/audio`,
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
            console.error('Error loading chunk:', error);
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
                            console.warn("Error removing source buffer", error);
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
                    console.error("Failed to play audio", error)
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

    // Changes the like button image
    function toggleLikedTrack() {
        isTrackLiked = !isTrackLiked
    }

</script>

<div class="w-72 h-auto bg-gray-700 rounded rounded-xl flex flex-col justify-center mb-4 mx-1" onmouseenter={onLikeAndDownload} onfocus={onLikeAndDownload} onmouseleave={offLikeAndDownload} role="button" tabindex="0">
    <!-- Picture section -->
    <div class="h-72 w-72 flex flex-row items-center justify-center">
        <div 
            id={trackId}
            onmouseover={hoverTrackImage} 
            onfocus={hoverTrackImage}
            onmouseleave={leaveHoverTrackImage}
            onclick={playPauseAudio} 
            onkeydown={(e) => {if (e.key === "Enter" || e.key === " ") playPauseAudio}} 
            role="button" 
            tabindex="0" 
            class="h-64 w-64 bg-slate-400 flex flex-row items-center rounded rounded-xl justify-center"
        >
            {#if coverURL}
                <img class="h-64 w-64 absolute" src={coverURL} alt="cover art">
            {/if}
            
            <!-- Absolutely disgusting conditional for whether to show play or pause button -->
            {#if isHovered}
                {#if $isPlaying}
                    {#if (trackId === $currentTrackId)}
                        <img class="h-20 w-20 absolute" src="pause.png" alt="Pause button" />
                    {:else}
                        <img class="h-20 w-20 absolute" src="play.png" alt="Play button" />
                    {/if}
                {:else}
                    <img class="h-20 w-20 absolute" src="play.png" alt="Play button" />
                {/if}
            {:else}
                {#if ($isPlaying && (trackId === $currentTrackId))}
                    <img class="h-20 w-20 absolute" src="pause.png" alt="Pause button" />
                {/if}
            {/if}
            
        </div>
    </div>

    <!-- Section below the picture -->  
    <div class="w-72 h-24 bg-gray-700 rounded rounded-xl px-4">
        <div class="flex flex-row w-full">
            <div class={`${isExpanded ? 'w-40' : 'w-64'}`}>
                <p class="hover:underline truncate">{trackName}</p>
                <p class="pb-2 truncate">{trackName}</p>
            </div>
            {#if isExpanded}
                <div class="w-24 flex flex-row">
                    <button class="w-12 h-12 flex flex-row items-center justify-center" onclick={toggleLikedTrack}>
                        {#if isTrackLiked}
                            <img class="h-8 w-8 hover:h-9 hover:w-9" src="heart-checked.png" alt="Like Button"/>                       
                        {:else}
                            <img class="h-8 w-8 hover:h-9 hover:w-9" src="heart-unchecked.png" alt="Like Button"/>
                        {/if}
                    </button>
                    <button class="w-12 h-12 flex flex-row items-center justify-center">
                        <img class="h-8 w-8 hover:h-9 hover:w-9" src="vinyl.png" alt="Layerr Button"/>
                    </button>
                </div>
            {/if}
        </div>
        <div>
            {#if previousTrackName}
                <button class="text-indigo-500 flex flex-row items-center">
                    <img class="w-6 h-6" src="vinyl.png" alt="Song samples" />
                    <p class="ml-2 hover:underline">[SAMPLE] {previousTrackName}</p>
                </button>
            {/if}
        </div>
    </div>
</div>
