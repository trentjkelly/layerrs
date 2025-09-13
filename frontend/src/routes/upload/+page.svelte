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
    let coverArtFiles = $state<FileList | null>(null);
    let artistLayerrs = $state<Array<string>>([]);
    let connections = $state<Array<string>>([]);
    let isConnectionsDropdownOpen = $state(false);
    let title = $state<string>('');
    let isUploaded = $state(false);
    let isDragOver = $state(false);

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

    function removeCoverArtFile() {
        coverArtFiles = null;
    }

    function toggleConnectionsDropdown() {
        isConnectionsDropdownOpen = !isConnectionsDropdownOpen;
    }

    function addConnection(trackId: string) {
        if (!connections.includes(trackId)) {
            connections = [...connections, trackId];
        }
        isConnectionsDropdownOpen = false;
    }

    function removeConnection(trackId: string) {
        connections = connections.filter(id => id !== trackId);
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
        if (imageFile) {
            const dt = new DataTransfer();
            dt.items.add(imageFile);
            coverArtFiles = dt.files;
        }
    }

    async function submitFile() {
        if (!audioFiles || !coverArtFiles) {
            logger.error("Missing audio or cover art files");
            return;
        }
        
        let audioFile = audioFiles[0];
        let coverArtFile = coverArtFiles[0];

        // Audio file validation: needs to be either wav or flac
        if (audioFile.type !== 'audio/wav' && audioFile.type !== 'audio/flac') {
            logger.error("audioFile is not a wav or flac file")
            return
        }

        // Cover art file validation: needs to be a png or jpg
        if (coverArtFile.type !== 'image/png' && coverArtFile.type !== 'image/jpeg') {
            logger.error("coverArtFile is not a png or jpg file")
            return
        }

        if (audioFile && coverArtFile) {
            logger.debug("audioFile and coverArtFile are valid")
            const form = new FormData();
            form.append('audioFile', audioFile)
            form.append('coverArtFile', coverArtFile)
            form.append('name', title)
            form.append('connections', JSON.stringify(connections))

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
            logger.error("audioFile and coverArtFile are not valid")
        }
    }

    function navigateHome() {
        goto('/')
    }

</script>

<main class={`transition-all duration-300 min-h-screen w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>

    <TopHeader pageName="Upload" pageIcon="/upload.png"></TopHeader>

    <section class="w-full flex flex-row justify-center pb-32">
        {#if $isLoggedIn}
            <div class="outline outline-gray-600 rounded-3xl w-2/3 max-w-4xl flex flex-col items-center p-8">
                {#if !isUploaded}
                <h2 class="mb-4 text-3xl font-bold text-white">Upload a Track</h2>
                
                <!-- Track Name Input -->
                <div class="w-full mb-6">
                    <h3 class="text-xl font-semibold text-white mb-3">Track Name</h3>
                    <input 
                        class="w-full px-4 py-3 rounded-lg bg-gray-700 text-white placeholder-gray-400 border border-gray-600 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" 
                        type="text" 
                        bind:value={title} 
                        placeholder="Enter track name..." 
                    />
                </div>

                <!-- Add Connections Section -->
                <div class="w-full mb-4">
                    <h3 class="text-xl font-semibold text-white mb-3">Add Layerrs</h3>
                    <div class="bg-gray-700 rounded-lg p-4 mb-4">
                        <p class="text-gray-200 text-sm leading-relaxed">
                            Layerrs are any other artist's tracks you've used for samples, vocals, sounds, remixes, covers, in this track.
                        </p>
                    </div>
                    
                    <!-- Selected Connections -->
                    {#if connections.length > 0}
                        <div class="mb-3 space-y-2">
                            {#each connections as trackId}
                                <div class="flex items-center justify-between px-3 bg-gray-600 rounded-lg">
                                    <span class="text-white text-md">{trackId}</span>
                                    <button 
                                        type="button"
                                        onclick={() => removeConnection(trackId)}
                                        class="text-red-400 hover:text-red-300 text-2xl font-bold"
                                    >
                                        ×
                                    </button>
                                </div>
                            {/each}
                        </div>
                    {/if}
                    
                    <div class="relative flex justify-center">
                        <button 
                            type="button"
                            onclick={toggleConnectionsDropdown}
                            class="px-6 py-3 rounded-full bg-violet-600 text-white hover:bg-violet-700 focus:outline-none focus:ring-2 focus:ring-violet-500/20 flex items-center justify-center transition-colors"
                        >
                            <span class="text-lg font-semibold">Add Layerrs +</span>
                        </button>
                        
                        {#if isConnectionsDropdownOpen}
                            <div class="absolute top-full left-0 right-0 mt-1 bg-gray-700 border border-gray-600 rounded-lg shadow-lg z-10 max-h-48 overflow-y-auto">
                                {#if artistLayerrs.length === 0}
                                    <div class="px-4 py-3 text-gray-400 text-center">
                                        No tracks available
                                    </div>
                                {:else}
                                    {#each artistLayerrs as trackId}
                                        <button 
                                            type="button"
                                            onclick={() => addConnection(trackId)}
                                            class="w-full px-4 py-3 text-left text-white hover:bg-gray-600 border-b border-gray-600 last:border-b-0"
                                        >
                                            {trackId}
                                        </button>
                                    {/each}
                                {/if}
                            </div>
                        {/if}
                    </div>
                </div>
                
                <!-- Audio Upload Box -->
                <div class="w-full mb-4">
                    <h3 class="text-xl font-semibold text-white mb-3">Audio File</h3>
                    <label for="audio" class="block">
                        <div 
                            role="button"
                            tabindex="0"
                            class="w-full h-48 border-2 border-dashed border-gray-400 rounded-xl flex flex-col items-center justify-center transition-all duration-200 cursor-pointer hover:border-indigo-400 hover:bg-gray-700/50 {isDragOver && !audioFiles ? 'border-indigo-500 bg-indigo-500/20' : ''}"
                            ondragover={handleDragOver}
                            ondragleave={handleDragLeave}
                            ondrop={handleDrop}
                        >
                            {#if !audioFiles}
                                <div class="text-center">
                                    <p class="text-lg text-gray-300">Drop your audio file here</p>
                                    <p class="text-sm text-gray-300">or click to browse</p>
                                    <p class="text-sm text-gray-400 mt-6">Only FLAC and WAV files are supported</p>
                                </div>
                            {:else}
                                <div class="text-center w-full">
                                    <div class="flex items-center justify-center space-x-2">
                                        <span class="text-green-400">✓</span>
                                        <span class="text-gray-300">{audioFiles[0].name}</span>
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

                <!-- Cover Art Upload Box -->
                <div class="w-full mb-4">
                    <h3 class="text-xl font-semibold text-white mb-3">Cover Art</h3>
                    <label for="coverArt" class="block">
                        <div 
                            role="button"
                            tabindex="0"
                            class="w-full h-48 border-2 border-dashed border-gray-400 rounded-xl flex flex-col items-center justify-center transition-all duration-200 cursor-pointer hover:border-indigo-400 hover:bg-gray-700/50 {isDragOver && !coverArtFiles ? 'border-indigo-500 bg-indigo-500/20' : ''}"
                            ondragover={handleDragOver}
                            ondragleave={handleDragLeave}
                            ondrop={handleDrop}
                        >
                            {#if !coverArtFiles}
                                <div class="text-center">
                                    <p class="text-lg text-gray-300">Drop your cover art here</p>
                                    <p class="text-sm text-gray-300">or click to browse</p>
                                    <p class="text-sm text-gray-400 mt-6">Only PNG and JPG files are supported</p>
                                </div>
                            {:else}
                                <div class="text-center w-full">
                                    <div class="flex items-center justify-center space-x-2">
                                        <span class="text-green-400">✓</span>
                                        <span class="text-gray-300">{coverArtFiles[0].name}</span>
                                        <button 
                                            onclick={removeCoverArtFile}
                                            class="ml-2 px-2 py-1 text-xs bg-red-500 hover:bg-red-600 rounded text-white"
                                        >
                                            Remove
                                        </button>
                                    </div>
                                </div>
                            {/if}
                        </div>
                    </label>
                    
                    <!-- Hidden cover art file input -->
                    <input id="coverArt" class="hidden" type="file" accept="image/*" bind:files={coverArtFiles} />
                </div>
                
                <button 
                    class="mt-8 px-8 py-4 bg-indigo-600 hover:bg-indigo-700 rounded-full text-white font-semibold text-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed" 
                    onclick={submitFile}
                    disabled={!audioFiles || !coverArtFiles || !title}
                >
                    Upload Track
                </button>
            {:else}
                <div class="w-full h-full flex flex-col items-center justify-center">
                    <h2 class="mb-8 text-3xl font-bold text-white">Track successfully uploaded!</h2>
                    <button class="bg-indigo-600 hover:bg-indigo-700 mb-12 px-8 py-4 rounded-full text-white font-semibold transition-colors" onclick={navigateHome}>
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

