<script lang="ts">
    import TopHeader from "../../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../../stores/player";
    import NewTrackCard from "../../../components/NewTrackCard.svelte";

    import { page } from '$app/state';
    import { urlBase } from "../../../stores/environment";
    import { jwt } from "../../../stores/auth";

    let slug = page.params.slug;
    let creditAgreement = $state(false);
    let noStealingAgreement = $state(false);
    let banAgreement = $state(false);
    let isSubmitting = $state(false);
    let downloaded = $state(false);
    
    async function handleSubmit(event: Event) {
        event.preventDefault();
        
        if (!creditAgreement || !noStealingAgreement || !banAgreement) {
            alert('Please check all agreement boxes before proceeding.');
            return;
        }
        
        isSubmitting = true;
        await handleDownload();
        isSubmitting = false;
        console.log("DONE")
    }

    async function handleDownload() {
        const response = await fetch(`${$urlBase}/api/track/${slug}/download`, {
            headers: {
                'Authorization': `Bearer ${$jwt}`
            }
        });

        if (!response.ok) {
            alert('Failed to download track');
            return;
        }

        const res = await response.json();

        if (res.url) {
            console.log(res.url);
            window.open(res.url);
            downloaded = true;

        } else {
            alert('Failed to download track');
        }
    };

</script>

<main class={`transition-all duration-300 h-full w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>
    <TopHeader pageName="" pageIcon=""></TopHeader>
    
    <section class="w-full flex flex-wrap justify-around pb-24">
        <NewTrackCard trackId={slug}></NewTrackCard>

        {#if !downloaded}
            <form class="w-3/4 max-w-[1200px] py-8 px-6 bg-zinc-800 rounded-lg border-zinc-700 mt-8" onsubmit={handleSubmit}>
                <div class="mb-6">
                    <h2 class="text-3xl font-bold text-white mb-4 text-center">Flip This Track</h2>
                    <div class="bg-zinc-700 rounded-lg p-4 mb-4">
                        <p class="text-zinc-200 text-sm leading-relaxed">
                            <strong class="text-white">Note:</strong> This track will be added to "Your Layerrs" and should be given proper credit when uploading any track that uses any part of this file.
                        </p>
                    </div>
                    <p class="text-zinc-300 text-base leading-relaxed mb-6">
                        Before proceeding, please read and agree to these terms:
                    </p>
                </div>

                <div class="mb-6 ml-2 space-y-4">
                    <div class="flex items-start space-x-3">
                        <input 
                            type="checkbox" 
                            id="credit-checkbox"
                            bind:checked={creditAgreement}
                            class="mt-1 w-5 h-5 text-violet-600 bg-zinc-700 border-zinc-600 rounded focus:ring-violet-500 hover:cursor-pointer focus:ring-2"
                        />
                        <label for="credit-checkbox" class="text-zinc-300 text-base leading-relaxed cursor-pointer">
                            When uploading a track that uses any part of this one, I will <strong class="text-white">give proper credit in the upload track form</strong>.
                        </label>
                    </div>

                    <div class="flex items-start space-x-3">
                        <input 
                            type="checkbox" 
                            id="no-stealing-checkbox"
                            bind:checked={noStealingAgreement}
                            class="mt-1 w-5 h-5 text-violet-600 bg-zinc-700 border-zinc-600 rounded focus:ring-violet-500 hover:cursor-pointer focus:ring-2"
                        />
                        <label for="no-stealing-checkbox" class="text-zinc-300 text-base leading-relaxed cursor-pointer">
                            I will <strong class="text-white">not steal the artist's work</strong> or claim it as my own.
                        </label>
                    </div>

                    <div class="flex items-start space-x-3">
                        <input 
                            type="checkbox" 
                            id="ban-checkbox"
                            bind:checked={banAgreement}
                            class="mt-1 w-5 h-5 text-violet-600 bg-zinc-700 border-zinc-600 rounded focus:ring-violet-500 hover:cursor-pointer focus:ring-2"
                        />
                        <label for="ban-checkbox" class="text-zinc-300 text-base leading-relaxed cursor-pointer">
                            I understand that <strong class="text-white">failure to abide by these terms may result in a permanent ban</strong>.
                        </label>
                    </div>
                </div>

                <div class="flex justify-center">
                    <button 
                        type="submit" 
                        disabled={!creditAgreement || !noStealingAgreement || !banAgreement || isSubmitting}
                        class="px-8 py-2 bg-violet-600 text-white rounded-lg hover:bg-violet-700 disabled:bg-zinc-600 disabled:cursor-not-allowed transition-all duration-300 flex items-center space-x-2"
                    >
                        {#if isSubmitting}
                            <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            <span>Downloading...</span>
                        {:else}
                            <span>Download Track</span>
                        {/if}
                    </button>
                </div>
            </form>
        {:else}
            <div class="w-3/4 max-w-[1200px] py-8 px-6 bg-zinc-800 rounded-lg border-zinc-700 mt-8">
                <div class="mb-6">
                    <h2 class="text-3xl font-bold text-white mb-4 text-center">Successfully Downloaded!</h2>
                    <p class="text-zinc-300 text-base leading-relaxed mb-6 text-center">
                        This track has been added to "Your Layerrs" and should be given proper credit when uploading any track that uses any part of this file. You can give credit and create a connection to this track when using the upload track form.
                    </p>
                </div>
            </div>
        {/if}
    </section>
</main>