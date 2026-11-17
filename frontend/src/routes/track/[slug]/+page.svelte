<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/state';
    import TopHeader from '../../../components/TopHeader.svelte';
    import ProjectDetail from '../../../components/ProjectDetail.svelte';
    import { getUrlBase } from '../../../stores/environment';

    type Detail = { name: string; description: string; artistName: string; parentId: number | null; parentName: string | null; tracks: { id: number; name: string; role: string; artistName: string; url: string }[] };
    let project = $state<Detail | null>(null);
    let error = $state('');
    onMount(async () => {
        try {
            const response = await fetch(`${getUrlBase()}/api/projects/${page.params.slug}/`);
            if (!response.ok) throw new Error('Unable to load this project.');
            project = await response.json();
            if (project) project.tracks = project.tracks.map(track => ({ ...track, url: track.url.startsWith('/') ? `${getUrlBase()}${track.url}` : track.url }));
        } catch (e) { error = e instanceof Error ? e.message : 'Unable to load this project.'; }
    });
</script>
<main class="min-h-full w-full bg-zinc-950">
    <TopHeader pageName="Project" pageIcon="" />
    <div class="mx-auto w-full max-w-4xl px-4 pb-32 pt-8">
        {#if project}<ProjectDetail {...project} />{:else if error}<p class="text-white" role="alert">{error}</p>{:else}<p class="text-zinc-300">Loading project…</p>{/if}
    </div>
</main>
