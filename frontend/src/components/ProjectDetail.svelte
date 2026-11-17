<script lang="ts">
    type AudioItem = { id: number; name: string; role: string; artistName: string; url: string };
    let { name, description, artistName, tracks, parentId = null, parentName = null }: {
        name: string; description: string; artistName: string; tracks: AudioItem[];
        parentId?: number | null; parentName?: string | null;
    } = $props();
</script>

<article class="w-full border border-zinc-700 bg-zinc-900 p-6 text-white sm:p-8">
    <p class="text-sm uppercase tracking-wide text-zinc-400">Project</p>
    <h1 class="mt-1 break-words text-3xl font-bold">{name}</h1>
    <p class="mt-2 text-zinc-300">by {artistName}</p>
    {#if parentId && parentName}<p class="mt-2 text-sm text-zinc-300">Built on <a class="text-violet-300 underline" href={`/track/${parentId}`}>{parentName}</a></p>{/if}
    {#if description}<p class="mt-5 whitespace-pre-wrap break-words text-zinc-200">{description}</p>{/if}
    <div class="mt-6 space-y-4">
        {#each tracks as track, index (`${track.role}-${track.id}-${index}`)}
            <section class="border border-zinc-700 bg-zinc-950 p-4">
                <p class="text-xs uppercase tracking-wide text-zinc-400">{track.role === 'master' ? 'Master' : 'Stem'}</p>
                <h2 class="mt-1 break-words font-semibold">{track.name}</h2>
                {#if track.role === 'stem' && track.artistName}<p class="text-sm text-zinc-400">Source: {track.artistName}</p>{/if}
                {#if track.url}<audio class="mt-3 w-full" controls preload="none" src={track.url} aria-label={`Play ${track.name}`}></audio>{:else}<p class="mt-3 text-sm text-zinc-400">Audio unavailable</p>{/if}
            </section>
        {/each}
    </div>
</article>
