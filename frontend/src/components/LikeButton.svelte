<script>
    import { jwt } from "../stores/auth";
    import { goto } from "$app/navigation";
    import { urlBase } from "../stores/environment";
    import { logger } from "../modules/lib/logger";

    let { trackId, numLikes, isLiked = false } = $props()
    let isTrackLiked = $state(isLiked)

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
            const res = await fetch(`${$urlBase}/api/likes`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${$jwt}`
                },
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
            const res = await fetch(`${$urlBase}/api/likes?${params}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${$jwt}`
                }
            })

            if (res.status == 401) {
                goto('/login')
            }

        } catch (error) {
            logger.error("Could not unlike track")
        }
    }

</script>

<button class="px-2 py-1 ml-2 flex flex-row items-center justify-center hover:bg-white transition-all duration-300 text-white hover:text-violet-500" onclick={toggleLikedTrack}>
    {#if isTrackLiked}
        <p>LIKED</p>                   
    {:else}
        <p>LIKE</p>
    {/if}

    {#if numLikes > 0}
        <p class="text-sm ml-2">{numLikes}</p>
    {/if}
</button>