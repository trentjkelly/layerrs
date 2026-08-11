<script lang="ts">
    import { goto } from "$app/navigation";
    import TopHeader from "../../components/TopHeader.svelte";
    import LayerrsTrackCard from "../../components/LayerrsTrackCard.svelte";
    import UserMenu from "../../components/UserMenu.svelte";
    import { isLoggedIn, authInitialized } from "../../stores/auth";
    import { fetchWithAuth } from "../../modules/lib/fetch";
    import { isSidebarOpen } from "../../stores/player";
    import { urlBase } from "../../stores/environment";
    import { logger } from "../../modules/lib/logger";
    import { onMount } from "svelte";
    import { get } from "svelte/store";
    import { usernameStore, emailStore, portraitUrlStore, loadProfile} from "../../stores/profile";

    type LayerrTrack = {
        id: number;
        description: string;
        artistName: string;
        artistPortraitUrl: string;
    };

    let audioFiles = $state<FileList | null>(null);
    let artistLayerrs = $state<Array<LayerrTrack>>([]);
    let layerrs = $state<Array<number>>([]);
    let description = $state<string>('');
    let priceInCents = $state<number>(0);
    let isUploaded = $state(false);
    let isDragOver = $state(false);
    let isLoading = $state(false);
    let username = $state('');

    let dollars = $state(0);
    let isPageLoaded = $state(false);
    
    $effect(() => {
        if ($authInitialized && !$isLoggedIn) {
            goto('/login');
        }
    });

    function display(d: number) {
        return d.toLocaleString() + '.00';
    }

    function onKeydown(e: KeyboardEvent) {
        if (e.key === 'Backspace') {
        e.preventDefault();
        dollars = Math.floor(dollars / 10);
        }
    }

    function onInput(e: any) {
        const digit = (e.data || '').replace(/\D/g, '');
        if (!digit) return;
        const next = dollars * 10 + parseInt(digit, 10);
        if (next <= 9999999) dollars = next;
    }

    function onPaste(e: ClipboardEvent) {
        e.preventDefault();
        if (e.clipboardData) {
            const pasted = e.clipboardData.getData('text');
            const digits = pasted.replace(/\D/g, '').slice(0, 7);
            if (digits) dollars = parseInt(digits, 10);
        } else {
            console.error("Clipboard data is null");
            dollars = 0;
        }
    }

    onMount(async () => {
        if ($isLoggedIn) {
            await getArtistLayerrs();
            await getArtistUsername();
        }
    })

    async function getArtistUsername() {
        username = get(usernameStore);
        if(username != ''){
            return;
        }
        usernameStore.set('');
        await loadProfile();
        username = get(usernameStore);
        if (username == '') {
            alert('You must set a username before continuing.');
            goto('/profile');
        }
    }

    async function getArtistLayerrs() {
        const response = await fetchWithAuth(`${$urlBase}/api/layerrs`);
        if (!response.ok) {
            throw new Error("Failed to get artist layerrs");
        }
        const layerrsData: Array<LayerrTrack> = await response.json();
        artistLayerrs = layerrsData ?? [];
    }

    function removeAudioFile() {
        audioFiles = null;
    }

    function addlayerr(trackId: number) {
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

        if (description.length < 3 || description.length > 100) {
            logger.error("Description must be between 3 and 100 characters")
            return
        }

        if (audioFile) {
            logger.debug("audioFile is valid")
            const form = new FormData();
            form.append('audioFile', audioFile)
            form.append('description', description)
            form.append('layerrIDs', JSON.stringify(layerrs))

            const res = await fetchWithAuth(`${$urlBase}/api/track/`, {
                method: "POST",
                body: form
            });
            if (res.status == 201) {
                const trackData = await res.json();
                const trackId = trackData.id;

                if (priceInCents >= 0) {
                    const priceInCentsValue = priceInCents * 100;
                    await fetchWithAuth(`${$urlBase}/api/track/${trackId}/price`, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({ priceInCents: priceInCentsValue })
                    });
                }

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

    $effect(() => {
        if ($authInitialized && $isLoggedIn) {
            isPageLoaded = true;
        }
    });

</script>

{#if isPageLoaded}
    <main class={`transition-all duration-300 h-screen overflow-y-auto w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-950`}>

        <TopHeader pageName="Upload" pageIcon="/upload.png"></TopHeader>

        <UserMenu username={$usernameStore} email={$emailStore} portraitUrl={$portraitUrlStore} />

        <section class="w-full flex flex-row justify-center pb-32">
            {#if $isLoggedIn}
                <div class="bg-zinc-900 border border-zinc-700 w-2/3 max-w-4xl flex flex-col items-center p-8">
                    {#if !isUploaded}
                    <h2 class="mb-4 text-3xl font-bold text-white">Upload a Track</h2>

                    <!-- Audio Upload Box -->
                    <div class="w-full mb-4">
                        <h3 class="text-xl font-semibold text-white mb-1">Audio File</h3>
                        <label for="audio" class="block">
                            <div 
                                role="button"
                                tabindex="0"
                                class="w-full h-48 border-2 border-dashed border-zinc-700 bg-zinc-800 flex flex-col items-center justify-center transition-all duration-200 cursor-pointer hover:bg-zinc-700 {isDragOver && !audioFiles ? 'border-focus bg-primary-muted' : ''}"
                                ondragover={handleDragOver}
                                ondragleave={handleDragLeave}
                                ondrop={handleDrop}
                            >
                                {#if !audioFiles}
                                    <div class="text-center">
                                        <p class="text-lg text-white">Drop your audio file here</p>
                                        <p class="text-sm text-white">or click to browse</p>
                                        <p class="text-sm text-white mt-6">Only FLAC and WAV files are supported</p>
                                    </div>
                                {:else}
                                    <div class="text-center w-full">
                                        <div class="flex items-center justify-center space-x-2">
                                            <span class="text-success">✓</span>
                                            <span class="text-white">{audioFiles[0].name}</span>
                                            <button 
                                                onclick={removeAudioFile}
                                                class="ml-2 px-2 py-1 text-xs bg-zinc-600 hover:bg-zinc-700 text-white cursor-pointer"
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
                        <h3 class="text-xl font-semibold text-white mb-1">Track Name / Description</h3>
                        <input
                            class="w-full px-2 py-2 bg-zinc-800 text-white placeholder-white border border-zinc-700 focus:border-focus focus:outline-hidden"
                            type="text"
                            bind:value={description}
                            placeholder="Give your track a short name or description..."
                            maxlength={100}
                        />
                        <p class="text-sm mt-1 {description.length === 0 ? 'text-white' : (description.length < 3 || description.length > 100 ? 'text-danger' : 'text-white')}">
                            {description.length}/100 characters (minimum 3)
                        </p>
                    </div>

                    <!-- Add layerrs Section -->
                    <div class="w-full mb-4">
                        <h3 class="text-xl font-semibold text-white mb-1">Add Layerrs</h3>
                        <p class="text-sm text-white mt-1 mb-3">Layerrs are any other artist's tracks you've used for samples, vocals, sounds, remixes, covers, in this track.</p>
                        
                        <div class="h-48 overflow-y-auto bg-zinc-900 border border-zinc-700">
                            {#if artistLayerrs.length === 0}
                                <div class="h-full flex items-center justify-center px-4 py-3 text-white text-center">
                                    You haven't downloaded any tracks yet
                                </div>
                            {:else}
                                {#each artistLayerrs as track}
                                    <LayerrsTrackCard
                                        {track}
                                        isSelected={layerrs.includes(track.id)}
                                        ontoggle={addlayerr}
                                    />
                                {/each}
                            {/if}
                        </div>
                    </div>

                    <!-- Price Input -->
                    <div class="w-full mb-4">
                        <h3 class="text-xl font-semibold text-white mb-1">Price (optional)</h3>
                        <p class="text-sm text-white mt-1 mb-3">Set a price in USD to sell your track. Leave at $0 for free.</p>
                        <div class="flex items-center gap-2">
                            <span class="text-white text-lg">$</span>
                                <input
                                class="bg-zinc-800 text-white border border-zinc-700 focus:border-focus focus:outline-hidden px-2 py-2"
                                type="text"
                                inputmode="numeric"
                                value={display(dollars)}
                                onkeydown={onKeydown}
                                oninput={onInput}
                                onpaste={onPaste}
                                />
                            <span class="text-white">USD</span>
                        </div>
                    </div>

                    <button
                        class="mt-8 px-8 py-4 bg-zinc-600 hover:bg-zinc-700 text-white font-semibold text-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3 cursor-pointer"
                        onclick={submitFile}
                        disabled={!audioFiles || description.length < 3 || description.length > 100 || isLoading}
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
                        <button class="bg-zinc-600 hover:bg-zinc-700 mb-2 px-8 py-4 text-white font-semibold transition-colors cursor-pointer" onclick={navigateHome}>
                            Return Home
                        </button>
                    </div>
                {/if}
            </div>
            {/if}
        </section>
    </main>
{/if}
