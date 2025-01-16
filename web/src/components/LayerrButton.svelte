<script>
    import { goto } from "$app/navigation";

    let { trackId, numLayerrs } = $props();

    function navigateLayerr() {
        goto(`/layerr/${trackId}`)
    }

</script>

<button class="w-12 h-12 flex flex-col items-center justify-center" onclick={navigateLayerr}>
    <img class="h-8 w-8 hover:h-9 hover:w-9" src="/vinyl.png" alt="Layerr Button"/>
    <p class="text-sm">{numLayerrs}</p>
</button>