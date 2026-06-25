<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { urlBase } from "../../../../stores/environment";
    import { isSidebarOpen } from "../../../../stores/player";
    import TopHeader from "../../../../components/TopHeader.svelte";
    import type { Page as PageModel, PageSubmission } from "../../../../models/types";
    import { getPage, getSubmissions, approveSubmission } from "../../../../modules/requests/page-requests";

    let pageId = $derived($page.params.pageId ?? "");
    let pageData: PageModel | null = $state(null);
    let submissions: PageSubmission[] = $state([]);
    let isLoading = $state(true);
    let errorMessage = $state("");

    let approvingId: number | null = $state(null);
    let approvalNote = $state("");

    async function loadData() {
        isLoading = true;
        errorMessage = "";

        const p = await getPage($urlBase, pageId);
        if (!p) {
            errorMessage = "Page not found";
            isLoading = false;
            return;
        }
        if (!p.isEditor) {
            goto(`/pages/${pageId}`);
            return;
        }

        pageData = p;
        submissions = await getSubmissions($urlBase, pageId);
        isLoading = false;
    }

    onMount(loadData);

    function openApproveModal(submissionId: number) {
        approvingId = submissionId;
        approvalNote = "";
    }

    function closeApproveModal() {
        approvingId = null;
        approvalNote = "";
    }

    async function handleApprove() {
        if (approvingId == null || !pageData) return;
        const pageTrack = await approveSubmission($urlBase, pageId, approvingId, approvalNote);
        if (pageTrack) {
            submissions = submissions.filter(s => s.id !== approvingId);
            closeApproveModal();
        }
    }

    const activeSubmission = $derived(submissions.find(s => s.id === approvingId));
</script>

<main class={`transition-all duration-300 min-h-screen w-full overflow-y-auto ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-olive-400`}>
    <TopHeader pageName="Submissions" pageIcon="" />

    <section class="p-8 max-w-3xl mx-auto">
        {#if isLoading}
            <p class="text-center">Loading...</p>
        {:else if errorMessage}
            <p class="text-center text-red-400">{errorMessage}</p>
        {:else}
            <div class="flex items-center justify-between mb-6">
                <h1 class="text-2xl font-bold">Submissions for {pageData?.name}</h1>
                <a href={`/pages/${pageId}`} class="text-blue-400 hover:underline">Back to Page</a>
            </div>

            {#if submissions.length === 0}
                <p class="text-zinc-300">No submissions yet.</p>
            {:else}
                <div class="flex flex-col gap-4">
                    {#each submissions as submission (submission.id)}
                        <div class="p-4 bg-zinc-800 rounded flex flex-col gap-2">
                            <div class="flex flex-row justify-between items-start">
                                <div>
                                    {#if submission.track}
                                        <p class="font-semibold">{submission.track.description || "Untitled track"}</p>
                                        <p class="text-sm text-zinc-300">by {submission.track.artistName || "Unknown artist"}</p>
                                    {:else}
                                        <p class="font-semibold">Track {submission.trackId}</p>
                                    {/if}
                                    {#if submission.submitter}
                                        <p class="text-xs text-zinc-400 mt-1">Submitted by {submission.submitter.username}</p>
                                    {/if}
                                </div>
                                <button
                                    onclick={() => openApproveModal(submission.id)}
                                    class="px-3 py-1 bg-green-600 hover:bg-green-700 rounded text-sm"
                                >
                                    Add to feed
                                </button>
                            </div>
                            {#if submission.note}
                                <p class="text-sm text-zinc-300 italic">“{submission.note}”</p>
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
        {/if}
    </section>

    {#if approvingId !== null && activeSubmission}
        <div class="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
            <div class="bg-zinc-900 p-6 rounded max-w-md w-full">
                <h2 class="text-xl font-bold mb-4">Add to Page feed?</h2>
                <p class="mb-2">
                    <span class="font-semibold">{activeSubmission.track?.description || "Untitled track"}</span>
                    {#if activeSubmission.track?.artistName}
                        <span class="text-zinc-300"> by {activeSubmission.track.artistName}</span>
                    {/if}
                </p>
                <label for="approval-note" class="block text-sm font-semibold mt-4 mb-1">Editor note (optional)</label>
                <textarea
                    id="approval-note"
                    bind:value={approvalNote}
                    rows="3"
                    class="w-full p-2 rounded bg-zinc-800 text-white"
                ></textarea>
                <div class="flex gap-3 mt-4 justify-end">
                    <button onclick={closeApproveModal} class="px-4 py-1 bg-zinc-700 rounded">Cancel</button>
                    <button onclick={handleApprove} class="px-4 py-1 bg-green-600 hover:bg-green-700 rounded">Confirm</button>
                </div>
            </div>
        </div>
    {/if}
</main>
