<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { urlBase } from "../stores/environment";
    import { getTrackPages } from "../modules/requests/page-requests";
    import type { TrackPagesResponse } from "../models/types";

    let { trackId }: { trackId: number } = $props();

    let data: TrackPagesResponse | null = $state(null);

    onMount(async () => {
        data = await getTrackPages($urlBase, String(trackId));
    });

    function formatFollowers(count: number): string {
        if (count >= 1_000_000) return (count / 1_000_000).toFixed(1) + "M";
        if (count >= 1_000) return (count / 1_000).toFixed(1) + "k";
        return String(count);
    }

    function navigate() {
        goto(`/tracks/${trackId}/pages`);
    }
</script>

{#if data?.topPage}
    <button
        type="button"
        onclick={navigate}
        class="text-left text-sm text-blue-300 hover:text-blue-200 hover:underline cursor-pointer"
    >
        {data.topPage.name} · {formatFollowers(data.topPage.followerCount || 0)}
        {#if data.otherCount > 0}
            <span class="text-zinc-400">+{data.otherCount} other page{data.otherCount === 1 ? "" : "s"}</span>
        {/if}
    </button>
{/if}
