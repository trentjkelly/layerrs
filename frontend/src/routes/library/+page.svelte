<script lang="ts">
    import TopHeader from "../../components/TopHeader.svelte";
    import LibraryTrackCard from "../../components/LibraryTrackCard.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { isLoggedIn, authInitialized } from "../../stores/auth";
    import { fetchWithAuth } from "../../modules/lib/fetch";
    import { urlBase } from "../../stores/environment";
    import type { TrackInfo } from "../../models/types";

    $effect(() => {
        if ($authInitialized && !$isLoggedIn) {
            goto('/login');
        }
    });

    let layerrTracks: TrackInfo[] = $state([]);
    let likedTracks: TrackInfo[] = $state([]);

    async function fetchLayerrs() {
        const response = await fetchWithAuth(`${$urlBase}/api/layerrs/`, { method: 'GET' });
        if (!response.ok) return;
        const data = await response.json();
        layerrTracks = (data as TrackInfo[]) ?? [];
    }

    async function fetchLikes() {
        const response = await fetchWithAuth(`${$urlBase}/api/recommendations/library/likes`, { method: 'GET' });
        if (!response.ok) return;
        const data = await response.json();
        likedTracks = (data as TrackInfo[]) ?? [];
    }

    onMount(async () => {
        await Promise.all([fetchLayerrs(), fetchLikes()]);
    });
</script>

<main class={`transition-all duration-300 h-screen overflow-y-auto overflow-x-hidden w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>
    <TopHeader pageName="Your Library" pageIcon="/vinyl.png"></TopHeader>

    <section class="w-full h-full flex flex-row justify-center pb-24">
        <div class="w-1/2">
            <h2 class="text-white text-2xl font-bold pl-8 pt-6 pb-2">My Layerrs</h2>
            <section class="mx-8 bg-zinc-700 border border-zinc-600 rounded-lg">
                {#if layerrTracks.length === 0}
                    <div class="flex items-center justify-center px-4 py-3 text-zinc-400">
                        You haven't downloaded any tracks yet
                    </div>
                {:else}
                    {#each layerrTracks as track}
                        <LibraryTrackCard {track} />
                    {/each}
                {/if}
            </section>
        
            <h2 class="text-white text-2xl font-bold pl-8 pt-6 pb-2">My Likes</h2>
            <section class="mx-8 bg-zinc-700 border border-zinc-600 rounded-lg">
                {#if likedTracks.length === 0}
                    <div class="flex items-center justify-center px-4 py-3 text-zinc-400">
                        You haven't liked any tracks yet
                    </div>
                {:else}
                    {#each likedTracks as track}
                        <LibraryTrackCard {track} />
                    {/each}
                {/if}
            </section>
        </div> 
    </section>
</main>
