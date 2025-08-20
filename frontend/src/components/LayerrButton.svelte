<script>
    import { goto } from "$app/navigation";

    let { trackId, numLayerrs } = $props();

    function navigateLayerr() {
        goto(`/layerrs/${trackId}`)
    }

</script>

<button class="w-12 h-12 flex flex-row items-center justify-center ml-4" onclick={navigateLayerr}>
    <img class="h-6 w-6 hover:h-7 hover:w-7" src="/vinyl.png" alt="Layerr Button"/>
</button>