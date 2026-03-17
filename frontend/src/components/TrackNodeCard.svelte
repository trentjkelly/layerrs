<script lang="ts">
    import { onMount } from 'svelte';

    type ColorName = 'red' | 'orange' | 'yellow' | 'green' | 'blue' | 'violet';

    const colorClasses: Record<ColorName, {
        bg: string;
        border: string;
        border700: string;
        text: string;
        badgeBg: string;
        badgeBorder: string;
        gradientFrom: string;
        shadowHex: string;
    }> = {
        red:    { bg: 'bg-red-500',    border: 'border-red-500',    border700: 'border-red-700',    text: 'text-red-400',    badgeBg: 'bg-red-500/20',    badgeBorder: 'border-red-500/40',    gradientFrom: 'from-red-500/25',    shadowHex: '#ef4444' },
        orange: { bg: 'bg-orange-500', border: 'border-orange-500', border700: 'border-orange-700', text: 'text-orange-400', badgeBg: 'bg-orange-500/20', badgeBorder: 'border-orange-500/40', gradientFrom: 'from-orange-500/25', shadowHex: '#f97316' },
        yellow: { bg: 'bg-yellow-500', border: 'border-yellow-500', border700: 'border-yellow-700', text: 'text-yellow-400', badgeBg: 'bg-yellow-500/20', badgeBorder: 'border-yellow-500/40', gradientFrom: 'from-yellow-500/25', shadowHex: '#eab308' },
        green:  { bg: 'bg-green-500',  border: 'border-green-500',  border700: 'border-green-700',  text: 'text-green-400',  badgeBg: 'bg-green-500/20',  badgeBorder: 'border-green-500/40',  gradientFrom: 'from-green-500/25',  shadowHex: '#22c55e' },
        blue:   { bg: 'bg-blue-500',   border: 'border-blue-500',   border700: 'border-blue-700',   text: 'text-blue-400',   badgeBg: 'bg-blue-500/20',   badgeBorder: 'border-blue-500/40',   gradientFrom: 'from-blue-500/25',   shadowHex: '#3b82f6' },
        violet: { bg: 'bg-violet-500', border: 'border-violet-500', border700: 'border-violet-700', text: 'text-violet-400', badgeBg: 'bg-violet-500/20', badgeBorder: 'border-violet-500/40', gradientFrom: 'from-violet-500/25', shadowHex: '#8b5cf6' },
    };

    type TrackNodeData = {
        id: number;
        artistName: string;
        description: string;
        artistPortraitUrl: string;
        duration: number;
        likes: number;
        color: string;
    };

    let { track, isActive = false, onclick = () => {} }: {
        track: TrackNodeData;
        isActive?: boolean;
        onclick?: () => void;
    } = $props();

    let hovered = $state(false);

    onMount(() => {
        console.log('[TrackNodeCard] track data:', track);
    });

    const colors = $derived(colorClasses[(track.color as ColorName) ?? 'violet']);
</script>

<div class="w-48 h-24 bg-zinc-800 rounded-2xl border-2 {colors.border700}">
    <div class="flex flex-row">
        {#if track.artistPortraitUrl}
            <img src={track.artistPortraitUrl} alt={track.artistName} class="w-7 h-7 rounded-lg object-cover ml-2" />
        {/if}
        <p class="ml-2">{track.artistName}</p>
    </div>
    <p>"{track.description}"</p>
</div>
