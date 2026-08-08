<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { isLoggedIn } from "../../../stores/auth";
    import { urlBase } from "../../../stores/environment";
    import { isSidebarOpen } from "../../../stores/player";
    import TopHeader from "../../../components/TopHeader.svelte";
    import PageTrackRow from "../../../components/PageTrackRow.svelte";
    import type { Page as PageModel, PageTrack } from "../../../models/types";
    import {
        getPage,
        updatePage,
        deletePage,
        followPage,
        unfollowPage,
        getPageFeed,
        removeTrackFromPage,
        submitTrack
    } from "../../../modules/requests/page-requests";

    let pageId = $derived($page.params.pageId ?? "");
    let pageData: PageModel | null = $state(null);
    let feed: PageTrack[] = $state([]);
    let isLoading = $state(true);
    let errorMessage = $state("");

    let isEditing = $state(false);
    let editName = $state("");
    let editDescription = $state("");

    let recommendTrackId = $state("");
    let recommendNote = $state("");
    let recommendError = $state("");

    async function loadData() {
        isLoading = true;
        errorMessage = "";

        const [p, f] = await Promise.all([
            getPage($urlBase, pageId),
            getPageFeed($urlBase, pageId)
        ]);

        if (p) {
            pageData = p;
            editName = p.name;
            editDescription = p.description;
        } else {
            errorMessage = "Page not found";
        }

        feed = f;
        isLoading = false;
    }

    onMount(loadData);

    async function handleUpdate() {
        if (!pageData) return;
        const updated = await updatePage($urlBase, pageId, editName.trim(), editDescription.trim());
        if (updated) {
            pageData = updated;
            isEditing = false;
        }
    }

    async function handleDelete() {
        if (!confirm("Are you sure you want to delete this Page?")) return;
        const ok = await deletePage($urlBase, pageId);
        if (ok) {
            goto("/");
        }
    }

    async function toggleFollow() {
        if (!pageData) return;
        const ok = pageData.isFollowing
            ? await unfollowPage($urlBase, pageId)
            : await followPage($urlBase, pageId);
        if (ok) {
            pageData.isFollowing = !pageData.isFollowing;
            pageData.followerCount = (pageData.followerCount || 0) + (pageData.isFollowing ? 1 : -1);
        }
    }

    async function handleRemoveTrack(trackId: number) {
        const ok = await removeTrackFromPage($urlBase, pageId, trackId);
        if (ok) {
            feed = feed.filter(item => item.trackId !== trackId);
        }
    }

    async function handleRecommend() {
        recommendError = "";
        const trackId = parseInt(recommendTrackId, 10);
        if (isNaN(trackId) || trackId <= 0) {
            recommendError = "Enter a valid track ID";
            return;
        }
        const submission = await submitTrack($urlBase, pageId, trackId, recommendNote);
        if (submission) {
            recommendTrackId = "";
            recommendNote = "";
            alert("Track recommended");
        } else {
            recommendError = "Failed to recommend track. Make sure you follow this Page and haven't already submitted it.";
        }
    }

    function getEditorName(page: PageModel | null): string {
        if (!page?.editorId) return "Anonymous";
        return page.editorName || "Anonymous";
    }
</script>

<main class={`transition-all duration-300 min-h-screen w-full overflow-y-auto ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-olive-400`}>
    <TopHeader pageName="Page" pageIcon="" />

    <section class="p-8 max-w-3xl mx-auto">
        {#if isLoading}
            <p class="text-center">Loading...</p>
        {:else if errorMessage}
            <p class="text-center text-red-400">{errorMessage}</p>
        {:else if pageData}
            <div class="mb-6">
                {#if isEditing}
                    <input
                        bind:value={editName}
                        class="w-full text-3xl font-bold mb-2 p-2 rounded bg-zinc-800 text-white"
                    />
                    <textarea
                        bind:value={editDescription}
                        rows="3"
                        class="w-full p-2 rounded bg-zinc-800 text-white"
                    ></textarea>
                    <div class="flex gap-2 mt-2">
                        <button onclick={handleUpdate} class="px-4 py-1 bg-blue-600 rounded">Save</button>
                        <button onclick={() => isEditing = false} class="px-4 py-1 bg-zinc-600 rounded">Cancel</button>
                    </div>
                {:else}
                    <h1 class="text-3xl font-bold">{pageData.name}</h1>
                    <p class="text-zinc-200 mt-2">{pageData.description || "No description"}</p>
                {/if}

                <div class="flex flex-row items-center gap-4 mt-4 text-sm text-zinc-300">
                    <span>{pageData.followerCount || 0} followers</span>
                    <span>•</span>
                    <span>Curated by {getEditorName(pageData)}</span>
                </div>

                <div class="flex flex-row items-center gap-3 mt-4">
                    {#if $isLoggedIn}
                        {#if pageData.isEditor}
                            <button onclick={() => isEditing = !isEditing} class="px-4 py-1 bg-zinc-700 hover:bg-zinc-600 rounded">
                                {isEditing ? "Cancel" : "Edit"}
                            </button>
                            <button onclick={handleDelete} class="px-4 py-1 bg-red-700 hover:bg-red-800 rounded">
                                Delete
                            </button>
                            <a href={`/pages/${pageId}/submissions`} class="px-4 py-1 bg-zinc-700 hover:bg-zinc-600 rounded">
                                Submissions
                            </a>
                        {:else}
                            <button onclick={toggleFollow} class="px-4 py-1 bg-blue-600 hover:bg-blue-700 rounded">
                                {pageData.isFollowing ? "Unfollow" : "Follow"}
                            </button>

                            {#if pageData.isFollowing}
                                <div class="flex flex-col gap-1 ml-4">
                                    <div class="flex flex-row gap-2">
                                        <input
                                            type="number"
                                            bind:value={recommendTrackId}
                                            placeholder="Track ID"
                                            class="px-2 py-1 rounded bg-zinc-800 text-white w-32"
                                        />
                                        <button onclick={handleRecommend} class="px-3 py-1 bg-green-600 hover:bg-green-700 rounded">
                                            Recommend track
                                        </button>
                                    </div>
                                    <textarea
                                        bind:value={recommendNote}
                                        placeholder="Optional note to the editor"
                                        rows="2"
                                        class="px-2 py-1 rounded bg-zinc-800 text-white w-80"
                                    ></textarea>
                                    {#if recommendError}
                                        <p class="text-red-400 text-xs">{recommendError}</p>
                                    {/if}
                                </div>
                            {/if}
                        {/if}
                    {/if}
                </div>
            </div>

            <h2 class="text-2xl font-bold mb-4">Feed</h2>
            {#if feed.length === 0}
                <p class="text-zinc-300">No tracks yet.</p>
            {:else}
                <div class="flex flex-col">
                    {#each feed as pageTrack (pageTrack.id)}
                        <PageTrackRow
                            {pageTrack}
                            isEditor={pageData.isEditor ?? false}
                            onRemove={() => handleRemoveTrack(pageTrack.trackId)}
                        />
                    {/each}
                </div>
            {/if}
        {/if}
    </section>
</main>
