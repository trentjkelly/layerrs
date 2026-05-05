<script lang="ts">
    import { onMount } from 'svelte';

    type TrackNodeData = {
        id: number;
        artistName: string;
        description: string;
        artistPortraitUrl: string;
        duration: number;
        likes: number;
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

</script>

<div class="w-48 h-24 bg-zinc-800 rounded-2xl border-2 border-zinc-700">
    <div class="flex flex-row">
        {#if track.artistPortraitUrl}
            <img src={track.artistPortraitUrl} alt={track.artistName} class="w-7 h-7 rounded-lg object-cover ml-2" />
        {/if}
        <p class="ml-2">{track.artistName}</p>
    </div>
    <p>"{track.description}"</p>
</div>
