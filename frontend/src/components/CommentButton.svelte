<script>
    import { goto } from '$app/navigation';

    let { trackId, numComments } = $props();

    function navigateComments() {
        goto(`/track/${trackId}/comments`)
    }
</script>

<div class="w-12 h-12 flex flex-row items-center justify-center ml-4">
    <img class="h-9 w-9 hover:h-10 hover:w-10 mb-1" src="/comment.png" alt="Comment Button"/>
    {#if numComments > 0}
        <p class="text-sm ml-2">{numComments}</p>
    {/if}
</div>