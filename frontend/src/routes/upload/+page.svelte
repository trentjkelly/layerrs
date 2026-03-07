<script lang="ts">
    import { goto } from "$app/navigation";
    import TopHeader from "../../components/TopHeader.svelte";
    import { isLoggedIn, jwt } from "../../stores/auth";
    import { isSidebarOpen } from "../../stores/player";
    import { handleEnvironment, urlBase } from "../../stores/environment";
    import { logger } from "../../modules/lib/logger";
    import { onMount } from "svelte";
    import LogInPopup from "../../components/LogInPopup.svelte";

    let audioFiles = $state<FileList | null>(null);
    let artistLayerrs = $state<Array<string>>([]);
    let layerrs = $state<Array<string>>([]);
    let description = $state<string>('');
    let isUploaded = $state(false);
    let isDragOver = $state(false);
    let isLoading = $state(false);

    onMount(async () => {
        await handleEnvironment();
        await getArtistLayerrs();
    })

    async function getArtistLayerrs() {
        const response = await fetch(`${$urlBase}/api/layerrs`, {
            headers: {
                'Authorization': `Bearer ${$jwt}`
            }
        });
        if (!response.ok) {
            throw new Error("Failed to get artist layerrs");
        }
        const layerrsData : Array<any> = await response.json();

        console.log(layerrsData);

        artistLayerrs = layerrsData.map(layerr => layerr.trackId);
    }

    function removeAudioFile() {
        audioFiles = null;
    }

    function addlayerr(trackId: string) {
        if (layerrs.includes(trackId)) {
            layerrs = layerrs.filter(id => id !== trackId);
        } else {
            layerrs = [...layerrs, trackId];
        }
    }

    function handleDragOver(event: DragEvent) {
        event.preventDefault();
        isDragOver = true;
    }

    function handleDragLeave(event: DragEvent) {
        event.preventDefault();
        isDragOver = false;
    }

    function handleDrop(event: DragEvent) {
        event.preventDefault();
        isDragOver = false;
        
        const files = Array.from(event.dataTransfer?.files || []);
        const audioFile = files.find(file => file.type.startsWith('audio/'));
        const imageFile = files.find(file => file.type.startsWith('image/'));
        
        if (audioFile) {
            audioFiles = new DataTransfer().files;
            const dt = new DataTransfer();
            dt.items.add(audioFile);
            audioFiles = dt.files;
        }
    }

    async function submitFile() {
        if (!audioFiles) {
            logger.error("Missing audio file");
            return;
        }
        isLoading = true;
        
        let audioFile = audioFiles[0];

        // Audio file validation: needs to be either wav or flac
        if (audioFile.type !== 'audio/wav' && audioFile.type !== 'audio/flac') {
            logger.error("audioFile is not a wav or flac file")
            return
        }

        if (description.length < 10 || description.length > 100) {
            logger.error("Description must be between 10 and 100 characters")
            return
        }

        if (audioFile) {
            logger.debug("audioFile is valid")
            const form = new FormData();
            form.append('audioFile', audioFile)
            form.append('description', description)
            form.append('layerrIDs', JSON.stringify(layerrs))

            const res = await fetch(`${$urlBase}/api/track/`, { 
                method: "POST", 
                headers: {
                    'Authorization': `Bearer ${$jwt}`
                },
                body: form
            });
            if (res.status == 201) {
                isUploaded = true
            }
        } else {
            logger.error("audioFile is not valid")
        }
        isLoading = false;
    }

    function navigateHome() {
        goto('/')
    }

</script>

<main class={`transition-all duration-300 min-h-screen w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>

    <TopHeader pageName="Upload" pageIcon="/upload.png"></TopHeader>

    <section class="w-full flex flex-row justify-center pb-32">
        {#if $isLoggedIn}
            <div class="outline outline-zinc-600 rounded-3xl w-2/3 max-w-4xl flex flex-col items-center p-8">
                {#if !isUploaded}
                <h2 class="mb-4 text-3xl font-bold text-white">Upload a Track</h2>

                <!-- Audio Upload Box -->
                <div class="w-full mb-4">
                    <h3 class="text-xl font-semibold text-white mb-1">Audio File</h3>
                    <label for="audio" class="block">
                        <div 
                            role="button"
                            tabindex="0"
                            class="w-full h-48 border-2 border-dashed border-zinc-400 rounded-xl flex flex-col items-center justify-center transition-all duration-200 cursor-pointer hover:border-indigo-400 hover:bg-zinc-700/50 {isDragOver && !audioFiles ? 'border-indigo-500 bg-indigo-500/20' : ''}"
                            ondragover={handleDragOver}
                            ondragleave={handleDragLeave}
                            ondrop={handleDrop}
                        >
                            {#if !audioFiles}
                                <div class="text-center">
                                    <p class="text-lg text-zinc-300">Drop your audio file here</p>
                                    <p class="text-sm text-zinc-300">or click to browse</p>
                                    <p class="text-sm text-zinc-400 mt-6">Only FLAC and WAV files are supported</p>
                                </div>
                            {:else}
                                <div class="text-center w-full">
                                    <div class="flex items-center justify-center space-x-2">
                                        <span class="text-green-400">✓</span>
                                        <span class="text-zinc-300">{audioFiles[0].name}</span>
                                        <button 
                                            onclick={removeAudioFile}
                                            class="ml-2 px-2 py-1 text-xs bg-red-500 hover:bg-red-600 rounded text-white"
                                        >
                                            Remove
                                        </button>
                                    </div>
                                </div>
                            {/if}
                        </div>
                    </label>
                    
                    <!-- Hidden audio file input -->
                    <input id="audio" class="hidden" type="file" accept="audio/*" bind:files={audioFiles} />
                </div>
                
                <!-- Description Input -->
                <div class="w-full mb-6">
                    <h3 class="text-xl font-semibold text-white mb-1">Description</h3>
                    <input
                        class="w-full px-2 py-2 rounded-lg bg-zinc-700 text-white placeholder-zinc-400 border border-zinc-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
                        type="text"
                        bind:value={description}
                        placeholder="Give your track a short description..."
                        maxlength={100}
                    />
                    <p class="text-sm mt-1 {description.length === 0 ? 'text-zinc-400' : (description.length < 10 || description.length > 100 ? 'text-red-400' : 'text-zinc-400')}">
                        {description.length}/100 characters (minimum 10)
                    </p>
                </div>

                <!-- Add layerrs Section -->
                <div class="w-full mb-4">
                    <h3 class="text-xl font-semibold text-white mb-1">Add Layerrs</h3>
                    <p class="text-sm text-zinc-400 mt-1 mb-3">Layerrs are any other artist's tracks you've used for samples, vocals, sounds, remixes, covers, in this track.</p>
                    
                    <div class="h-48 overflow-y-auto bg-zinc-700 border border-zinc-600 rounded-lg">
                        {#if artistLayerrs.length === 0}
                            <div class="h-full flex items-center justify-center px-4 py-3 text-zinc-400 text-center">
                                You haven't downloaded any tracks yet
                            </div>
                        {:else}
                            {#each artistLayerrs as trackId}
                                <button
                                    type="button"
                                    onclick={() => addlayerr(trackId)}
                                    class="w-full flex items-center justify-between px-4 py-3 text-left text-white border-b border-zinc-600 last:border-b-0 cursor-pointer transition-colors {layerrs.includes(trackId) ? 'bg-violet-900/30 hover:bg-violet-900/50' : 'hover:bg-zinc-600'}"
                                >
                                    <span>{trackId}</span>
                                    {#if layerrs.includes(trackId)}
                                        <span class="text-violet-400 font-bold">✓</span>
                                    {/if}
                                </button>
                            {/each}
                        {/if}
                    </div>
                </div>
                
                <button
                    class="mt-8 px-8 py-4 bg-indigo-600 hover:bg-indigo-700 rounded-full text-white font-semibold text-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3"
                    onclick={submitFile}
                    disabled={!audioFiles || description.length < 10 || description.length > 100 || isLoading}
                >
                    {#if isLoading}
                        <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                        </svg>
                        Uploading...
                    {:else}
                        Upload Track
                    {/if}
                </button>
            {:else}
                <div class="w-full h-full flex flex-col items-center justify-center">
                    <h2 class="mb-4 text-3xl font-bold text-white">Track successfully uploaded!</h2>
                    <button class="bg-indigo-600 hover:bg-indigo-700 mb-2 px-8 py-4 rounded-full text-white font-semibold transition-colors" onclick={navigateHome}>
                        Return Home
                    </button>
                </div>
            {/if}
        </div>
        {:else}
            <LogInPopup />
        {/if}
    </section>
</main>

