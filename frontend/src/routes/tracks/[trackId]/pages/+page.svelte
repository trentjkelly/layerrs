<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { isSidebarOpen } from "../../../../stores/player";
    import { urlBase } from "../../../../stores/environment";
    import TopHeader from "../../../../components/TopHeader.svelte";
    import type { PageWithFollowerCount } from "../../../../models/types";
    import { getTrackPages } from "../../../../modules/requests/page-requests";
    import { formatFollowers } from "../../../../modules/lib/format";

    let trackId = $derived($page.params.trackId ?? "");
    let pages: PageWithFollowerCount[] = $state([]);
    let topPage: PageWithFollowerCount | undefined = $state(undefined);
    let otherCount = $state(0);
    let isLoading = $state(true);
    let errorMessage = $state("");

    async function loadData() {
        isLoading = true;
        errorMessage = "";

        const response = await getTrackPages($urlBase, trackId);
        if (response) {
            pages = response.pages;
            topPage = response.topPage;
            otherCount = response.otherCount;
        } else {
            errorMessage = "Failed to load Pages for this track";
        }

        isLoading = false;
    }

    onMount(loadData);
</script>

<main class={`transition-all duration-300 min-h-screen w-full overflow-y-auto ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-olive-400`}>
    <TopHeader pageName="Track Pages" pageIcon="" />

    <section class="p-8 max-w-3xl mx-auto">
        <h1 class="text-2xl font-bold mb-6">Pages featuring this track</h1>

        {#if isLoading}
            <p class="text-center">Loading...</p>
        {:else if errorMessage}
            <p class="text-center text-red-400">{errorMessage}</p>
        {:else if pages.length === 0}
            <p class="text-zinc-300">This track hasn't been added to any Pages yet.</p>
        {:else}
            {#if topPage}
                <div class="mb-6 p-4 bg-zinc-800 rounded">
                    <p class="text-sm text-zinc-400 mb-1">Top Page</p>
                    <a href={`/pages/${topPage.id}`} class="text-xl font-bold hover:underline">
                        {topPage.name}
                    </a>
                    <p class="text-zinc-300 mt-1">{formatFollowers(topPage.followerCount || 0)} followers</p>
                    <p class="text-sm text-zinc-400 mt-1">Curated by {topPage.editorName || "Anonymous"}</p>
                </div>
            {/if}

            <div class="flex flex-col gap-3">
                {#each pages as p (p.id)}
                    <a href={`/pages/${p.id}`} class="p-4 bg-zinc-800/50 rounded hover:bg-zinc-800 transition-colors">
                        <div class="flex flex-row justify-between items-center">
                            <div>
                                <p class="font-semibold text-lg">{p.name}</p>
                                <p class="text-sm text-zinc-400">Curated by {p.editorName || "Anonymous"}</p>
                            </div>
                            <span class="text-sm text-zinc-300">{formatFollowers(p.followerCount || 0)} followers</span>
                        </div>
                    </a>
                {/each}
            </div>

            {#if otherCount > 0}
                <p class="text-sm text-zinc-400 mt-4">+{otherCount} other page{otherCount === 1 ? "" : "s"} include this track</p>
            {/if}
        {/if}
    </section>
</main>
