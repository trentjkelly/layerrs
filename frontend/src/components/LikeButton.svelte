<script lang="ts">
    import { untrack } from 'svelte';
    import { goto } from "$app/navigation";
    import { urlBase } from "../stores/environment";
    import { logger } from "../modules/lib/logger";
    import { fetchWithAuth } from "../modules/lib/fetch";

    let { trackId, numLikes, isLiked = false} = $props()

    let isTrackLiked = $state(untrack(() => isLiked))

    // Changes the like button image, numLikes, and requests backend to save a like
    async function toggleLikedTrack() {
        isTrackLiked = !isTrackLiked

        if (isTrackLiked) {
            numLikes++
            await sendLikeRequest()
        } else {
            numLikes--
            await sendUnlikeRequest()
        }
    }

    // Requests the backend to like a track for the given user
    async function sendLikeRequest() {
        const formData = new FormData();
        formData.append('trackId', trackId)

        try {
            const res = await fetchWithAuth(`${$urlBase}/api/likes`, {
                method: 'POST',
                body: formData
            })

            if (res.status == 401) {
                goto('/login')
            }

        } catch (error) {
            logger.error("Could not like track")
        }
    }
    
    // Requests the backend to unlike a track for the given user
    async function sendUnlikeRequest() {
        const params = new URLSearchParams({
            trackId: trackId
        })

        try {
            const res = await fetchWithAuth(`${$urlBase}/api/likes?${params}`, { method: 'DELETE' })

            if (res.status == 401) {
                goto('/login')
            }

        } catch (error) {
            logger.error("Could not unlike track")
        }
    }

</script>

<button class="mx-1 px-2 py-1 outline-1 hover:bg-olive-600 hover:cursor-pointer flex flex-row items-center" onclick={toggleLikedTrack}>
    {#if isTrackLiked}
        <img src="heart-checked.png" alt="Liked" class="w-4 h-4">
        <p class="ml-1">Liked</p>                   
    {:else}
        <img src="heart-unchecked.png" alt="Like" class="w-4 h-4">
        <p class="ml-1">Like</p>
    {/if}

    {#if numLikes > 0}
        <p class="text-sm ml-2">{numLikes}</p>
    {/if}
</button>