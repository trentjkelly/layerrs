<script>
    import { goto } from "$app/navigation";

    let { trackId, numLayerrs } = $props();

    function navigateLayerr() {
        goto(`/layerrs/${trackId}`)
    }

</script>

<button class="py-1 px-2 ml-4 flex flex-row items-center justify-center hover:bg-white transition-all duration-300 text-white hover:text-violet-500" onclick={navigateLayerr}>
    <p class="text-md">FLIP THIS TRACK</p>
</button>