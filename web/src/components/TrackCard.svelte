<script>
    import { onMount } from 'svelte';
    import { currentSong, isPlaying } from '../stores/player';

    // TODO: Delete these variables
    let sampleTrackName = "Whole Lotta Red Whole Lotta Red"
    let sampleArtistName = "Playboi Carti Playboi Carti Playboi Carti"
    let samplePreviousSong = "Whole lotta blue"
        
    // Inherits the trackId from the page
    let { trackId } = $props();
    let coverURL = $state();
    let audioURL = $state();
    let trackName = $state();

    let isExpanded = $state(false);
    let isTrackLiked = $state(false)

    // Gets the cover art for the track when the component is loaded
    onMount(async () => {
        await getTrackData()
        await getCover()
    })

    async function getTrackData() {
        try {
            const response = await fetch(`http://localhost:8080/api/track/${trackId}/data`, { method: "GET"});
            if (!response.ok) {
                throw new Error("Failed to get track data");
            }
            const trackData = await response.json();
            trackName = trackData.name
        } catch (error) {
            console.error("Error catching track data", error)
        }
    }

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

    // Used when someone hovers over the track image
    async function hoverTrackImage() {
        // Add a play button
        const container = document.getElementById('imageDiv');
        const overlay = document.createElement('img');

        if ($isPlaying) {
            overlay.src = '/pause.png'
        } else {
            overlay.src = '/play.png'
        }

        overlay.alt = 'Overlay Image';
        overlay.id = 'overlay'
        overlay.style.position = 'absolute';
        overlay.style.pointerEvents = 'none';
        overlay.height = 60;
        overlay.width = 60;

        if (container) {
            container.append(overlay)
        }
        
        // Get the track audio
        // await getAudio()
    }

    async function leaveHoverTrackImage() {
        // Remove a play button
        const overlay = document.getElementById('overlay');
        if (overlay) {
            overlay.remove()
        }
    }

    // Gets the audio for the track when someone hovers over the component
    async function getAudio() {
        try {
            const response = await fetch(`http://localhost:8080/api/track/${trackId}/audio`, { method: "GET"});
            if (!response.ok) {
                throw new Error("Failed to stream track audio");
            }
            const blob = await response.blob();
            audioURL = URL.createObjectURL(blob);

        } catch (error) {
            console.error("Error streaming track audio", error)
        }
    }

    // Plays the audio 
    function playPauseAudio() {
        currentSong.set(audioURL)
    }

    function onLikeAndDownload() {
        isExpanded = true
    }

    function offLikeAndDownload() {
        isExpanded = false
    }

    function toggleLikedTrack() {
        isTrackLiked = !isTrackLiked
    }

</script>

<div class="w-72 h-auto bg-gray-700 rounded rounded-xl flex flex-col justify-center mb-4 mx-1" onmouseenter={onLikeAndDownload} onfocus={onLikeAndDownload} onmouseleave={offLikeAndDownload} role="button" tabindex="0">
    <!-- Picture section -->
    <div class="h-72 w-72 flex flex-row items-center justify-center">
        <div 
            id="imageDiv"
            onmouseover={hoverTrackImage} 
            onmouseleave={leaveHoverTrackImage}
            onfocus={hoverTrackImage} 
            onclick={playPauseAudio} 
            onkeydown={(e) => {if (e.key === "ENTER" || e.key === " ") playPauseAudio}} 
            role="button" 
            tabindex="0" 
            class="h-64 w-64 bg-slate-400 flex flex-row items-center rounded rounded-xl justify-center"
        >
            {#if coverURL}
                <img class="h-64 w-64 absolute" src={coverURL} alt="cover art">
            {/if}
        </div>
    </div>

    <!-- Section below the picture -->  
    <div class="w-72 h-24 bg-gray-700 rounded rounded-xl px-4">
        <div class="flex flex-row w-full">
            <div class={`${isExpanded ? 'w-40' : 'w-64'}`}>
                <p class="hover:underline truncate">{sampleTrackName}</p>
                <p class="pb-2 truncate">{sampleArtistName}</p>
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
            <button class="text-indigo-500 flex flex-row items-center">
                <img class="w-6 h-6" src="vinyl.png" alt="Song samples" />
                <p class="ml-2 hover:underline">[SAMPLE] {samplePreviousSong}</p>
            </button>
        </div>
    </div>
</div>
