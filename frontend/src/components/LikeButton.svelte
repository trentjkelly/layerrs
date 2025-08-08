<script>
    import { isLoggedIn } from "../stores/auth";
    import { jwt } from "../stores/auth";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { urlBase } from "../stores/environment";
    import { logger } from "../lib/logger";
    
    let { trackId, numLikes } = $props()
    let isTrackLiked = $state(false)

    onMount(async () => {
        await getIsLiked()
    })

    // Checks if the track is liked when page is loaded
    async function getIsLiked() {
        // If user is logged in
        if ($isLoggedIn) {
            try {
                const params = new URLSearchParams({
                    trackId: trackId
                })

                const response = await fetch(`${$urlBase}/api/likes?${params}`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${$jwt}`
                    }
                })
                if (response.ok) {
                    const data = await response.json()
                    isTrackLiked = data.isLiked
                }
            } catch (error) {
                logger.error(error)
            }
        }
    }

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

<button class="w-12 h-12 flex flex-col items-center justify-center" onclick={toggleLikedTrack}>
    {#if isTrackLiked}
        <img class="h-8 w-8 hover:h-9 hover:w-9" src="/heart-checked.png" alt="Like Button"/>                       
    {:else}
        <img class="h-8 w-8 hover:h-9 hover:w-9" src="/heart-unchecked.png" alt="Like Button"/>
    {/if}
    <p class="text-sm">{numLikes}</p>
</button>