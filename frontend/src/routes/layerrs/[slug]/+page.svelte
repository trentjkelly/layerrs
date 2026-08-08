<script lang="ts">
    import { onMount } from 'svelte';
    import TopHeader from "../../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../../stores/player";
    import NewTrackCard from "../../../components/NewTrackCard.svelte";

    import { goto } from '$app/navigation';
    import { page } from '$app/state';
    import { getUrlBase } from "../../../stores/environment";
    import { urlBase } from "../../../stores/environment";
    import { isLoggedIn, authInitialized } from "../../../stores/auth";
    import type { TrackInfo } from "../../../models/types";
    import { getTrackTrackInfo } from "../../../modules/requests/track-requests";
    import { fetchWithAuth } from "../../../modules/lib/fetch";
    import UploadTrackCard from "../../../components/UploadTrackCard.svelte";

    let slug = page.params.slug;
    let track = $state<TrackInfo | undefined>(page.state.track);
    let isLoadingTrack = $state(false);

    onMount(async () => {
        if (!track) {
            isLoadingTrack = true;
            if (slug){
                track = await getTrackTrackInfo(getUrlBase(), slug) ?? undefined;
            } else {
                console.error("No slug provided")
            }
            isLoadingTrack = false;
        }
    });
    let termsAgreement = $state(false);
    let isSubmitting = $state(false);

    async function handleSubmit(event: Event) {
        event.preventDefault();

        if (!termsAgreement) {
            alert('Please check the agreement box before proceeding.');
            return;
        }

        isSubmitting = true;
        await handleDownload();
        isSubmitting = false;
        console.log("DONE")
    }

    async function handleDownload() {
        const response = await fetchWithAuth(`${$urlBase}/api/track/${slug}/download`);

        if (!response.ok) {
            alert('Failed to download track');
            return;
        }

        const res = await response.json();

        if (res.url) {
            console.log(res.url);
            window.open(res.url);
            goto('/upload');

        } else {
            alert('Failed to download track');
        }
    };

</script>

<main class={`transition-all duration-300 h-screen overflow-y-auto w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-olive-400`}>
    <TopHeader pageName="" pageIcon=""></TopHeader>

    <section class="w-full flex flex-col items-center pb-32">
        {#if $isLoggedIn}
            <form class="bg-olive-500 border border-white w-2/3 max-w-4xl flex flex-col items-center p-8 mt-8" onsubmit={handleSubmit}>
                <h2 class="mb-4 text-3xl font-bold text-white">Build on This Track</h2>

                <!-- Track Card -->
                {#if isLoadingTrack}
                    <div class="w-full mb-6 flex justify-center">
                        <p class="text-white">Loading track info...</p>
                    </div>
                {:else if track}
                    <div class="w-full mb-6 flex justify-center">
                        <UploadTrackCard track={{ id: track.id, description: track.description, artistName: track.artistName, artistPortraitUrl: track.artistPortraitUrl }} />
                    </div>
                {/if}

                <!-- Note Box -->
                <div class="w-full mb-6">
                    <div class="bg-olive-600 border border-white p-4">
                        <p class="text-white text-sm leading-relaxed">
                            <strong class="text-white">Note:</strong> This track will be added to "Your Layerrs" and should be given proper credit when uploading any track that uses any part of this file.
                        </p>
                    </div>
                </div>

                <!-- Agreement Section -->
                <div class="w-full mb-6">
                    <h3 class="text-xl font-semibold text-white mb-3">Before proceeding, please read and agree to these terms:</h3>

                    <div class="bg-olive-600 border border-white p-4">
                        <ul class="list-disc list-outside ml-6 space-y-2 text-white text-base leading-relaxed">
                            <li>
                                When uploading a track that uses any part of this one, I will <strong class="text-white">give proper credit in the upload track form</strong>.
                            </li>
                            <li>
                                I will <strong class="text-white">not steal the artist's work</strong> or claim it as my own.
                            </li>
                            <li>
                                I understand that <strong class="text-white">failure to abide by these terms may result in a permanent ban</strong>.
                            </li>
                        </ul>

                        <div class="flex items-start space-x-3 mt-4 pt-4 border-t border-white/30">
                            <input
                                type="checkbox"
                                id="terms-checkbox"
                                bind:checked={termsAgreement}
                                class="mt-1 w-5 h-5 bg-olive-600 border border-white hover:cursor-pointer flex-shrink-0"
                            />
                            <label for="terms-checkbox" class="text-white text-base leading-relaxed cursor-pointer">
                                I have read and agree to the terms above.
                            </label>
                        </div>
                    </div>
                </div>

                <button
                    type="submit"
                    disabled={!termsAgreement || isSubmitting}
                    class="mt-4 px-8 py-4 bg-olive-600 hover:bg-olive-700 text-white font-semibold text-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3 cursor-pointer"
                >
                    {#if isSubmitting}
                        <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                        </svg>
                        Downloading...
                    {:else}
                        Download Track
                    {/if}
                </button>
            </form>
        {/if}
    </section>
</main>