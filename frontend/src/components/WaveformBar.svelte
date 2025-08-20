<script>
    import { currentTime } from '../stores/player';

    let { height, timePerBar, index, cursorTime = 0, isCursorHovered = false } = $props();

    let thresholdTime = $derived(timePerBar * index);
    let isActive = $derived($currentTime > thresholdTime);
    let isCursorActive = $derived(isCursorHovered && cursorTime > thresholdTime);
    let isTimeActiveAfterCursor = $derived(isCursorHovered && isActive && thresholdTime > cursorTime);
    let shouldHighlight = $derived(isCursorActive || isTimeActiveAfterCursor || (!isCursorHovered && isActive));

</script>

{#if shouldHighlight}
    <div
        class="w-[2px] mr-[2px] rounded rounded-xl"
        class:bg-violet-500={isCursorActive}
        class:bg-violet-300={isTimeActiveAfterCursor || (!isCursorHovered && isActive)}
        style="height: {height}%"
    ></div>
{:else}
    <div
        class="w-[2px] bg-white mr-[2px] rounded rounded-xl"
        style="height: {height}%"
    ></div>
{/if}
