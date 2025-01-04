<script>
    import { onMount } from 'svelte';
    import { currentSong, isPlaying } from '../stores/player';

    // Inherits the trackId from the page
    let { trackId } = $props();
    let coverURL = $state();
    let audioURL = $state();
    let trackName = $state();

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

</script>

<div class="w-72 h-96 bg-slate-200 flex flex-col justify-center">
    <div 
        id="imageDiv"
        onmouseover={hoverTrackImage} 
        onmouseleave={leaveHoverTrackImage}
        onfocus={hoverTrackImage} 
        onclick={playPauseAudio} 
        onkeydown={(e) => {if (e.key === "ENTER" || e.key === " ") playPauseAudio}} 
        role="button" 
        tabindex="0" 
        class="h-72 w-72 bg-slate-400 flex flex-row items-center justify-center"
    >
        {#if coverURL}
            <img class="h-72 w-72 absolute" src={coverURL} alt="cover art">
        {/if}
    </div>
    <div class="w-72 h-24">
        <p>{trackName}</p>
    </div>
</div>