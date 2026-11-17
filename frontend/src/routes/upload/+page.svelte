<script lang="ts">
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";
    import { get } from "svelte/store";
    import ProjectDetail from "../../components/ProjectDetail.svelte";
    import TopHeader from "../../components/TopHeader.svelte";
    import LayerrsTrackCard from "../../components/LayerrsTrackCard.svelte";
    import { isLoggedIn, authInitialized } from "../../stores/auth";
    import { fetchWithAuth } from "../../modules/lib/fetch";
    import { urlBase } from "../../stores/environment";
    import { logger } from "../../modules/lib/logger";
    import { usernameStore, loadProfile } from "../../stores/profile";

    type LayerrTrack = { id: number; description: string; artistName: string; artistPortraitUrl: string };
    type UploadType = 'original' | 'addon' | null;

    let uploadType = $state<UploadType>(null);
    let masterFiles = $state<FileList | null>(null);
    let stemFiles = $state<FileList | null>(null);
    let artistLayerrs = $state<LayerrTrack[]>([]);
    let selectedProjectId = $state<number | null>(null);
	let sourceTrackIds = $state<number[]>([]);
    let description = $state('');
    let projectDescription = $state('');
    let stemNames = $state<string[]>([]);
    let sourceNames = $state<Record<number, string>>({});
    let reviewing = $state(false);
    let createdId = $state<number | null>(null);
    let previewTracks = $state<{id: number; name: string; role: string; artistName: string; url: string}[]>([]);
    let objectUrls: string[] = [];
    function clearPreview() { objectUrls.forEach(URL.revokeObjectURL); objectUrls = []; previewTracks = []; }
    function backToEdit() { clearPreview(); reviewing = false; }
    async function showPreview() {
        if (!canUpload || !masterFile) return;
        clearPreview();
        const makeURL = (file: File) => { const url = URL.createObjectURL(file); objectUrls.push(url); return url; };
        const items = [{ id: 0, name: description.trim(), role: 'master', artistName: $usernameStore, url: makeURL(masterFile) },
            ...Array.from(stemFiles ?? []).map((file, index) => ({ id: index + 1, name: stemNames[index].trim(), role: 'stem', artistName: $usernameStore, url: makeURL(file) }))];
        for (const id of sourceTrackIds) {
            const track = artistLayerrs.find(t => t.id === id);
            let url = '';
            try { const response = await fetch(`${$urlBase}/api/track/${id}/audio`); if (response.ok) url = (await response.json()).url; } catch { /* displayed as unavailable */ }
            items.push({ id, name: sourceNames[id].trim(), role: 'stem', artistName: track?.artistName ?? '', url });
        }
        previewTracks = items;
        reviewing = true;
    }
    onMount(() => () => clearPreview());
    let isUploaded = $state(false);
    let isMasterDragOver = $state(false);
    let isStemsDragOver = $state(false);
    let isLoading = $state(false);
    let isPageLoaded = $state(false);
    let masterError = $state('');
    let stemsError = $state('');
    let submissionError = $state('');

    let isDescriptionValid = $derived(description.trim().length >= 3 && description.trim().length <= 100);
    let masterFile = $derived(masterFiles?.[0] ?? null);
    let hasValidStems = $derived((stemFiles?.length ?? 0) + sourceTrackIds.length <= 5 && (!stemFiles || Array.from(stemFiles).every(isSupportedAudio)));
    let canUpload = $derived(uploadType !== null && (uploadType !== 'addon' || selectedProjectId !== null) && isDescriptionValid && masterFile !== null && isSupportedAudio(masterFile) && hasValidStems && stemNames.every(name => name.trim().length > 0 && name.trim().length <= 255) && sourceTrackIds.every(id => (sourceNames[id] ?? '').trim().length > 0 && sourceNames[id].trim().length <= 255) && !isLoading);

    $effect(() => { if ($authInitialized && !$isLoggedIn) goto('/login'); });
    $effect(() => { if ($authInitialized && $isLoggedIn) isPageLoaded = true; });

    function isSupportedAudio(file: File) {
        return /\.(wav|flac)$/i.test(file.name);
    }

    function asFileList(files: File[]) {
        const dataTransfer = new DataTransfer();
        files.forEach((file) => dataTransfer.items.add(file));
        return dataTransfer.files;
    }

    function chooseUploadType(type: Exclude<UploadType, null>) {
        uploadType = type;
        selectedProjectId = null;
		sourceTrackIds = [];
        description = '';
        projectDescription = '';
        stemNames = []; sourceNames = {};
        masterFiles = null;
        stemFiles = null;
        masterError = '';
        stemsError = '';
        submissionError = '';
    }

    async function getArtistUsername() {
        if (get(usernameStore) !== '') return;
        usernameStore.set('');
        await loadProfile();
        if (get(usernameStore) === '') {
            alert('You must set a username before continuing.');
            goto('/profile');
        }
    }

    async function getArtistLayerrs() {
        const response = await fetchWithAuth(`${$urlBase}/api/layerrs`);
        if (!response.ok) throw new Error('Failed to get your available projects');
        artistLayerrs = (await response.json()) ?? [];
    }

    let profileLoaded = false;
    $effect(() => {
        if ($authInitialized && $isLoggedIn && !profileLoaded) {
            profileLoaded = true;
            void Promise.all([getArtistLayerrs(), getArtistUsername()]).catch((error) => { submissionError = error instanceof Error ? error.message : 'Unable to load your projects.'; });
        }
    });

    function selectProject(trackId: number) {
        selectedProjectId = selectedProjectId === trackId ? null : trackId;
    }

	function toggleSourceTrack(trackId: number) {
		if (sourceTrackIds.includes(trackId)) {
			sourceTrackIds = sourceTrackIds.filter((id) => id !== trackId);
			stemsError = '';
			return;
		}
		if ((stemFiles?.length ?? 0) + sourceTrackIds.length >= 5) {
			stemsError = 'You can include up to 5 uploaded or referenced stems.';
			return;
		}
		sourceTrackIds = [...sourceTrackIds, trackId];
        sourceNames = { ...sourceNames, [trackId]: artistLayerrs.find(t => t.id === trackId)?.description ?? '' };
		stemsError = '';
	}

    function setMasterFile(files: File[]) {
        masterError = '';
        const audioFiles = files.filter(isSupportedAudio);
        if (audioFiles.length === 0) {
            masterFiles = null;
            masterError = 'Choose one WAV or FLAC master file.';
            return;
        }
        if (audioFiles.length > 1) masterError = 'Only one master file can be uploaded.';
        masterFiles = asFileList([audioFiles[0]]);
    }

    function setStemFiles(files: File[]) {
        stemsError = '';
        if (files.length + sourceTrackIds.length > 5) {
            stemFiles = null; stemNames = [];
            stemsError = 'You can upload up to 5 stems.';
        } else if (!files.every(isSupportedAudio)) {
            stemFiles = null; stemNames = [];
            stemsError = 'Stem files must be WAV or FLAC.';
        } else {
            stemFiles = files.length ? asFileList(files) : null;
            stemNames = files.map(file => file.name);
        }
    }

    function handleDrop(event: DragEvent, target: 'master' | 'stems') {
        event.preventDefault();
        isMasterDragOver = false;
        isStemsDragOver = false;
        const files = Array.from(event.dataTransfer?.files ?? []);
        target === 'master' ? setMasterFile(files) : setStemFiles(files);
    }

    function removeStemFile(index: number) {
        if (stemFiles) { const names = stemNames.filter((_, i) => i !== index); setStemFiles(Array.from(stemFiles).filter((_, i) => i !== index)); stemNames = names; }
    }

    async function submitFile() {
        submissionError = '';
        if (!canUpload || !masterFile) return;
        isLoading = true;
        try {
            const form = new FormData();
            form.append('audioFile', masterFile);
            form.append('name', description.trim());
            form.append('projectDescription', projectDescription.trim());
            form.append('stemNames', JSON.stringify([...stemNames.map(n => n.trim()), ...sourceTrackIds.map(id => sourceNames[id].trim())]));
            form.append('layerrIDs', JSON.stringify(selectedProjectId === null ? [] : [selectedProjectId]));
			form.append('sourceTrackIDs', JSON.stringify(sourceTrackIds));
            Array.from(stemFiles ?? []).forEach((stem) => form.append('stemFiles', stem));
            const response = await fetchWithAuth(`${$urlBase}/api/projects/`, { method: 'POST', body: form });
            if (!response.ok) throw new Error(await response.text() || 'The project could not be uploaded.');
            createdId = (await response.json()).id;
            clearPreview();
            isUploaded = true;
        } catch (error) {
            logger.error(error);
            submissionError = error instanceof Error ? error.message : 'The project could not be uploaded.';
        } finally {
            isLoading = false;
        }
    }
</script>

{#if isPageLoaded}
    <main class="min-h-full w-full bg-zinc-950">
        <TopHeader pageName="Upload" pageIcon="/upload.png" />
        <section class="w-full flex justify-center pb-32">
            {#if $isLoggedIn}
                <div class="bg-zinc-900 border border-zinc-700 w-2/3 max-w-4xl flex flex-col items-center p-8">
                    {#if !isUploaded && reviewing}
                        <h2 class="mb-6 text-3xl font-bold text-white">Preview your project</h2>
                        <ProjectDetail name={description.trim()} description={projectDescription.trim()} artistName={$usernameStore} tracks={previewTracks} parentId={selectedProjectId} parentName={artistLayerrs.find(t => t.id === selectedProjectId)?.description ?? null} />
                        {#if submissionError}<p class="mt-4 text-danger" role="alert">{submissionError}</p>{/if}
                        <div class="mt-6 flex gap-3"><button class="bg-zinc-700 px-6 py-3 text-white" onclick={backToEdit} disabled={isLoading}>Back to edit</button><button class="bg-zinc-600 px-6 py-3 font-semibold text-white disabled:opacity-50" onclick={submitFile} disabled={isLoading}>{isLoading ? 'Posting…' : 'Post project'}</button></div>
                    {:else if !isUploaded}
                        <h2 class="mb-6 text-3xl font-bold text-white">Upload a Project</h2>
                        <div class="w-full">
                            <h3 class="text-xl font-semibold text-white mb-3">What are you uploading?</h3>
                            <div class="grid gap-3 sm:grid-cols-2">
                                <button class={`border p-4 text-left transition-colors cursor-pointer ${uploadType === 'original' ? 'border-focus bg-primary-muted' : 'border-zinc-700 bg-zinc-800 hover:bg-zinc-700'}`} onclick={() => chooseUploadType('original')}>
                                    <span class="block text-lg font-semibold text-white">Original project</span>
                                    <span class="mt-1 block text-sm text-zinc-300">Start a new project from your own material.</span>
                                </button>
                                <button class={`border p-4 text-left transition-colors cursor-pointer ${uploadType === 'addon' ? 'border-focus bg-primary-muted' : 'border-zinc-700 bg-zinc-800 hover:bg-zinc-700'}`} onclick={() => chooseUploadType('addon')}>
                                    <span class="block text-lg font-semibold text-white">Add-on to another project</span>
                                    <span class="mt-1 block text-sm text-zinc-300">Build on a project you have downloaded.</span>
                                </button>
                            </div>
                        </div>

                        {#if uploadType === 'addon'}
                            <div class="w-full mt-6" in:fade={{ duration: 180 }}>
                                <h3 class="text-xl font-semibold text-white mb-1">Select the project you are adding to</h3>
                                <p class="text-sm text-zinc-300 mt-1 mb-3">Choose one project from your Layerrs.</p>
                                <div class="h-48 overflow-y-auto bg-zinc-950 border border-zinc-700">
                                    {#if artistLayerrs.length === 0}
                                        <div class="h-full flex items-center justify-center px-4 py-3 text-zinc-300 text-center">You have not downloaded any projects yet.</div>
                                    {:else}
                                        {#each artistLayerrs as track}
                                            <LayerrsTrackCard {track} isSelected={selectedProjectId === track.id} ontoggle={selectProject} />
                                        {/each}
                                    {/if}
                                </div>
                            </div>
                        {/if}

                        {#if uploadType === 'original' || selectedProjectId !== null}
                            <div class="w-full mt-6" in:fade={{ duration: 180 }}>
                                <h3 class="text-xl font-semibold text-white mb-1">Project name</h3>
                                <input class="w-full px-2 py-2 bg-zinc-800 text-white placeholder-zinc-400 border border-zinc-700 focus:border-focus focus:outline-hidden" type="text" bind:value={description} placeholder="Name your project..." maxlength={100} />
                                <p class={`text-sm mt-1 ${description.length === 0 || isDescriptionValid ? 'text-zinc-300' : 'text-danger'}`}>{description.length}/100 characters (minimum 3)</p>
                            <label class="mt-4 block text-xl font-semibold text-white" for="project-description">Description (optional)</label>
                                <textarea id="project-description" class="mt-1 w-full border border-zinc-700 bg-zinc-800 p-2 text-white" bind:value={projectDescription} maxlength={2000} rows="4" placeholder="Tell people about your project..."></textarea>
                                <p class="text-sm text-zinc-300">{projectDescription.length}/2000 characters</p>
                            </div>
                        {/if}

                        {#if isDescriptionValid && (uploadType === 'original' || selectedProjectId !== null)}
                            <div class="w-full mt-6" in:fade={{ duration: 180 }}>
                                <h3 class="text-xl font-semibold text-white mb-1">{uploadType === 'original' ? 'Master track' : 'Upload your version'}</h3>
                                <p class="text-sm text-zinc-300 mt-1 mb-3">Upload one WAV or FLAC master file.</p>
                                <label for="master-audio" class="block">
                                    <div role="button" tabindex="0" class={`w-full min-h-40 border-2 border-dashed bg-zinc-800 flex flex-col items-center justify-center transition-all duration-200 cursor-pointer hover:bg-zinc-700 ${isMasterDragOver && !masterFile ? 'border-focus bg-primary-muted' : 'border-zinc-700'}`} ondragover={(event) => { event.preventDefault(); isMasterDragOver = true; }} ondragleave={() => isMasterDragOver = false} ondrop={(event) => handleDrop(event, 'master')}>
                                        {#if !masterFile}
                                            <div class="text-center"><p class="text-lg text-white">Drop your master file here</p><p class="text-sm text-zinc-300">or click to browse</p></div>
                                        {:else}
                                            <div class="flex items-center gap-2 px-4 text-center"><span class="text-success">✓</span><span class="text-white break-all">{masterFile.name}</span><button type="button" onclick={() => { masterFiles = null; masterError = ''; }} class="px-2 py-1 text-xs bg-zinc-600 hover:bg-zinc-700 text-white cursor-pointer">Remove</button></div>
                                        {/if}
                                    </div>
                                </label>
                                <input id="master-audio" class="hidden" type="file" accept=".wav,.flac,audio/wav,audio/flac" onchange={(event) => setMasterFile(Array.from((event.currentTarget as HTMLInputElement).files ?? []))} />
                                {#if masterError}<p class="text-danger text-sm mt-2">{masterError}</p>{/if}
                            </div>
                        {/if}

                        {#if masterFile && isSupportedAudio(masterFile) && isDescriptionValid}
                            <div class="w-full mt-6" in:fade={{ duration: 180 }}>
                                <h3 class="text-xl font-semibold text-white mb-1">{uploadType === 'original' ? 'Stems (optional)' : 'Stems or your part solo’d (optional)'}</h3>
                                <p class="text-sm text-zinc-300 mt-1 mb-3">Upload up to 5 WAV or FLAC files.</p>
                                <label for="stem-audio" class="block">
                                    <div role="button" tabindex="0" class={`w-full min-h-32 border-2 border-dashed bg-zinc-800 flex flex-col items-center justify-center transition-all duration-200 cursor-pointer hover:bg-zinc-700 ${isStemsDragOver ? 'border-focus bg-primary-muted' : 'border-zinc-700'}`} ondragover={(event) => { event.preventDefault(); isStemsDragOver = true; }} ondragleave={() => isStemsDragOver = false} ondrop={(event) => handleDrop(event, 'stems')}>
                                        <div class="text-center px-4"><p class="text-lg text-white">Drop stem files here</p><p class="text-sm text-zinc-300">or click to browse</p></div>
                                    </div>
                                </label>
                                <input id="stem-audio" class="hidden" type="file" multiple accept=".wav,.flac,audio/wav,audio/flac" onchange={(event) => setStemFiles(Array.from((event.currentTarget as HTMLInputElement).files ?? []))} />
                                {#if stemFiles}
                                    <ul class="mt-3 space-y-2">
                                        {#each Array.from(stemFiles) as stem, index}
                                            <li class="flex items-center justify-between gap-3 bg-zinc-800 px-3 py-2 text-sm text-white"><label class="flex min-w-0 flex-1 flex-col gap-1"><span class="break-all text-zinc-400">{stem.name}</span><span class="sr-only">Name for {stem.name}</span><input class="w-full border border-zinc-600 bg-zinc-950 p-2 text-white" bind:value={stemNames[index]} maxlength="255" /></label><button type="button" onclick={() => removeStemFile(index)} class="shrink-0 px-2 py-1 text-xs bg-zinc-600 hover:bg-zinc-700 cursor-pointer">Remove</button></li>
                                        {/each}
                                    </ul>
                                {/if}
                                {#if stemsError}<p class="text-danger text-sm mt-2">{stemsError}</p>{/if}
                            </div>
                            <div class="w-full mt-6" in:fade={{ duration: 180 }}>
                                <h3 class="text-xl font-semibold text-white mb-1">Use downloaded tracks as stems (optional)</h3>
                                <p class="text-sm text-zinc-300 mt-1 mb-3">Select tracks you have downloaded. Uploaded and referenced stems share a 5-stem limit.</p>
                                <div class="h-48 overflow-y-auto bg-zinc-950 border border-zinc-700">
                                    {#if artistLayerrs.length === 0}
                                        <div class="h-full flex items-center justify-center px-4 py-3 text-zinc-300 text-center">You have not downloaded any projects yet.</div>
                                    {:else}
                                        {#each artistLayerrs as track}
                                            <LayerrsTrackCard {track} isSelected={sourceTrackIds.includes(track.id)} ontoggle={toggleSourceTrack} />
                                        {/each}
                                    {/if}
                                </div>
                            </div>
                            {#each sourceTrackIds as sourceId}
                                <label class="mt-3 block w-full text-sm text-zinc-300">Name for {artistLayerrs.find(t => t.id === sourceId)?.description}
                                    <input class="mt-1 w-full border border-zinc-700 bg-zinc-800 p-2 text-white" bind:value={sourceNames[sourceId]} maxlength="255" />
                                </label>
                            {/each}
                            <div class="w-full mt-8 flex flex-col items-center" in:fade={{ duration: 180 }}>
                                {#if submissionError}<p class="mb-3 text-danger" role="alert">{submissionError}</p>{/if}
                                <button class="px-8 py-4 bg-zinc-600 hover:bg-zinc-700 text-white font-semibold text-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3 cursor-pointer" onclick={showPreview} disabled={!canUpload}>
                                    {#if isLoading}
                                        <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" aria-hidden="true"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path></svg>
                                        Uploading...
                                    {:else}
                                        Preview Project
                                    {/if}
                                </button>
                            </div>
                        {/if}
                    {:else}
                        <div class="w-full min-h-80 flex flex-col items-center justify-center" in:fade={{ duration: 180 }}>
                            <h2 class="mb-4 text-3xl font-bold text-white">Project successfully uploaded!</h2>
                            <button class="bg-zinc-600 hover:bg-zinc-700 mb-2 px-8 py-4 text-white font-semibold cursor-pointer" onclick={() => goto(`/track/${createdId}`)}>View Project</button>
                            <button class="bg-zinc-600 hover:bg-zinc-700 mb-2 px-8 py-4 text-white font-semibold transition-colors cursor-pointer" onclick={() => goto('/')}>Return Home</button>
                        </div>
                    {/if}
                </div>
            {/if}
        </section>
    </main>
{/if}
