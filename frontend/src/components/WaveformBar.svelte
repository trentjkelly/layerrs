<script lang="ts">
    let { height, timePerBar, index, cursorTime = 0, isCursorHovered = false, playedPercentage = 0, trackDuration = 0 } = $props();


    let thresholdTime = $derived(timePerBar * index);
    let isActive = $derived(playedPercentage * trackDuration > thresholdTime);
    let isCursorActive = $derived(isCursorHovered && cursorTime > thresholdTime);
    let isTimeActiveAfterCursor = $derived(isCursorHovered && isActive && thresholdTime > cursorTime);
    let shouldHighlight = $derived(isCursorActive || isTimeActiveAfterCursor || (!isCursorHovered && isActive));

</script>

{#if shouldHighlight}
    <div
        class="w-[2px] shrink-0 mr-[2px] bg-black"
        style="height: {height}%"
    ></div>
{:else}
    <div
        class="w-[2px] shrink-0 bg-black mr-[2px]"
        style="height: {height}%"
    ></div>
{/if}
