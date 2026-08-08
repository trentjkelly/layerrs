<script lang="ts">
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { isLoggedIn } from "../../stores/auth";
    import { urlBase } from "../../stores/environment";
    import { artistIdStore } from "../../stores/profile";
    import TopHeader from "../../components/TopHeader.svelte";
    import { getArtistPages } from "../../modules/requests/page-requests";
    import type { Page } from "../../models/types";

    let pages: Page[] = $state([]);
    let isLoading = $state(true);

    async function loadPages() {
        const artistId = $artistIdStore;
        if (!artistId) {
            isLoading = false;
            return;
        }
        pages = await getArtistPages($urlBase, String(artistId));
        isLoading = false;
    }

    onMount(() => {
        if ($isLoggedIn) {
            loadPages();
        } else {
            isLoading = false;
        }
    });
</script>

<main class={`transition-all duration-300 min-h-screen w-full overflow-y-auto ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-olive-400`}>
    <TopHeader pageName="Pages" pageIcon="/vinyl.png" />

    <section class="p-8 max-w-3xl mx-auto">
        <div class="flex flex-row items-center justify-between mb-6">
            <h1 class="text-3xl font-bold">Your Pages</h1>
            {#if $isLoggedIn}
                <button
                    onclick={() => goto('/pages/new')}
                    class="px-6 py-2 bg-blue-600 hover:bg-blue-700 rounded font-semibold"
                >
                    Create New Page
                </button>
            {/if}
        </div>

        {#if !$isLoggedIn}
            <p class="text-center text-lg">You must be logged in to view your Pages.</p>
        {:else if isLoading}
            <p class="text-center">Loading...</p>
        {:else if pages.length === 0}
            <p class="text-zinc-300">You don't have any Pages yet.</p>
        {:else}
            <div class="flex flex-col gap-3">
                {#each pages as page (page.id)}
                    <a
                        href={`/pages/${page.id}`}
                        class="block p-4 bg-olive-500 hover:bg-olive-600 rounded transition-colors"
                    >
                        <h2 class="text-xl font-semibold">{page.name}</h2>
                        {#if page.description}
                            <p class="text-zinc-200 mt-1 line-clamp-2">{page.description}</p>
                        {/if}
                        <p class="text-sm text-zinc-300 mt-2">{page.followerCount ?? 0} followers</p>
                    </a>
                {/each}
            </div>
        {/if}
    </section>
</main>
