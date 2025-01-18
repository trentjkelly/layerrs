<script>
    import { goto } from "$app/navigation";

    let { parentTrackId, parentTrackName, parentTrackArist } = $props();

    function navigateParentTrack() {
        goto(`/track/${parentTrackId}`)
    }

</script>

<button class="text-indigo-500 flex flex-row items-center mt-2" onclick={navigateParentTrack}>
    <img class="w-6 h-6 hover:w-7 hover:h-7" src="/vinyl-purple-3.png" alt="Song samples" />
    <p class="ml-2 hover:underline text-md">{parentTrackName} - @{parentTrackArist}</p>
</button>