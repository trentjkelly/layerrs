<script>
    import { currentSong, isPlaying } from "../stores/player";
</script>

<div class="w-screen h-20 bg-slate-700 fixed bottom-0 flex flex-row items-center justify-center">
    <button class="h-12 w-12 rounded-full bg-indigo-400">
        {#if isPlaying}
            <img class="scale-50 h-12 w-12" alt="Pause button" src="pause.png">
        {:else}
            <img class="scale-50 h-12 w-12" alt="Play button" src="play.png">
        {/if}
    </button>
</div>