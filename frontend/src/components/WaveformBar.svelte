<script lang="ts">
    import { currentTime } from '../stores/player';

    type ColorName = 'red' | 'orange' | 'yellow' | 'green' | 'blue' | 'violet';

    const colorClasses: Record<ColorName, { bg500: string; bg300: string }> = {
        red:    { bg500: 'bg-red-500',    bg300: 'bg-red-300' },
        orange: { bg500: 'bg-orange-500', bg300: 'bg-orange-300' },
        yellow: { bg500: 'bg-yellow-500', bg300: 'bg-yellow-300' },
        green:  { bg500: 'bg-green-500',  bg300: 'bg-green-300' },
        blue:   { bg500: 'bg-blue-500',   bg300: 'bg-blue-300' },
        violet: { bg500: 'bg-violet-500', bg300: 'bg-violet-300' },
    };

    let { height, timePerBar, index, cursorTime = 0, isCursorHovered = false, color = 'violet' } = $props();

    const colors = $derived(colorClasses[(color as ColorName) ?? 'violet']);

    let thresholdTime = $derived(timePerBar * index);
    let isActive = $derived($currentTime > thresholdTime);
    let isCursorActive = $derived(isCursorHovered && cursorTime > thresholdTime);
    let isTimeActiveAfterCursor = $derived(isCursorHovered && isActive && thresholdTime > cursorTime);
    let shouldHighlight = $derived(isCursorActive || isTimeActiveAfterCursor || (!isCursorHovered && isActive));

</script>

{#if shouldHighlight}
    <div
        class="w-[2px] mr-[2px] rounded-sm rounded-xl {isCursorActive ? colors.bg500 : colors.bg300}"
        style="height: {height}%"
    ></div>
{:else}
    <div
        class="w-[2px] bg-white mr-[2px] rounded-sm rounded-xl"
        style="height: {height}%"
    ></div>
{/if}
