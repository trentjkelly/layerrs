<script>

    let { trackId } = $props();
    
    /**
     * @type {string | undefined}
     */
    let coverURL = $state();

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

</script>

<div class="w-full h-96 bg-slate-200 flex justify-center">
    <div class="h-72 w-72 bg-slate-400">
        {#if coverURL }
            <img class="h-full w-full" alt="cover art">
        {:else}
            <p>Loading...</p>
        {/if}
    </div>
</div>