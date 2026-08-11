<script lang="ts">
    import type { PageTrack } from "../models/types";

    let { pageTrack, isEditor, onRemove }: { pageTrack: PageTrack; isEditor: boolean; onRemove?: () => void } = $props();

    const addedAt = $derived(pageTrack.addedAt ? new Date(pageTrack.addedAt).toLocaleDateString() : "");
</script>

<div class="w-full py-3 border-b border-zinc-700 flex flex-col gap-2">
    <div class="flex flex-row items-center justify-between">
        <div class="flex flex-col">
            {#if pageTrack.track}
                <a href={`/track/${pageTrack.track.id}`} class="text-lg font-semibold hover:underline">
                    {pageTrack.track.description || "Untitled track"}
                </a>
                <a href={`/artist/${pageTrack.track.artistId}`} class="text-sm text-zinc-300 hover:underline">
                    {pageTrack.track.artistName || "Unknown artist"}
                </a>
            {:else}
                <span class="text-lg font-semibold">Track {pageTrack.trackId}</span>
            {/if}
        </div>

        <div class="flex flex-row items-center gap-4">
            <span class="text-sm text-zinc-400">Added {addedAt}</span>
            {#if isEditor && onRemove}
                <button
                    onclick={onRemove}
                    class="text-sm px-3 py-1 bg-danger-solid hover:bg-danger-solid-hover"
                >
                    Remove
                </button>
            {/if}
        </div>
    </div>

    {#if pageTrack.notes && pageTrack.notes.length > 0}
        <div class="bg-zinc-800 p-3 mt-1">
            {#each pageTrack.notes as note}
                <p class="text-sm text-zinc-200 italic">“{note.note}”</p>
            {/each}
        </div>
    {/if}
</div>
