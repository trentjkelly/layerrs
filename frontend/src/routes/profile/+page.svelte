<script lang="ts">
    import { onMount } from "svelte";
    import { get } from "svelte/store";
    import { page } from "$app/state";
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { urlBase } from "../../stores/environment";
    import { fetchWithAuth } from "../../modules/lib/fetch";
    import { logger } from "../../modules/lib/logger";
    import { usernameStore, emailStore, bioStore, validStripeSellerStore, portraitUrlStore, loadProfile } from "../../stores/profile";

    let username = $state('');
    let bio = $state('');
    let validStripeSeller = $state(false);

    let portraitUrl = $state<string | null>(null);
    let photoFile = $state<File | null>(null);
    let isSaving = $state(false);
    let isOnboarding = $state(false);
    let saveError = $state<string | null>(null);
    let saveSuccess = $state(false);
    let usernameError = $state<string | null>(null);

    onMount(async () => {
        await loadProfile();
        username = get(usernameStore);
        bio = get(bioStore);
        portraitUrl = get(portraitUrlStore) || null;
        validStripeSeller = get(validStripeSellerStore);

        const onboarding = page.url.searchParams.get('onboarding');
        if (onboarding === 'complete') {
            onboardingComplete();
        } else if (onboarding === 'refresh') {
            await onboardSeller();
        }
    });

    function onboardingComplete() {
        console.log('stripe account is complete');
    }

    function handlePhotoChange(event: Event) {
        const input = event.target as HTMLInputElement;
        const file = input.files?.[0] ?? null;
        if (!file) return;

        const img = new Image();
        const objectUrl = URL.createObjectURL(file);

        img.onload = () => {
            URL.revokeObjectURL(objectUrl);

            const size = Math.min(img.width, img.height);
            const canvas = document.createElement('canvas');
            canvas.width = size;
            canvas.height = size;

            const ctx = canvas.getContext('2d')!;
            ctx.drawImage(
                img,
                (img.width - size) / 2,
                (img.height - size) / 2,
                size, size,
                0, 0, size, size
            );

            const outputType = file.type === 'image/png' ? 'image/png' : 'image/jpeg';
            canvas.toBlob((blob) => {
                if (!blob) return;
                photoFile = new File([blob], file.name, { type: outputType });
                portraitUrl = URL.createObjectURL(photoFile);
            }, outputType, 0.9);
        };

        img.src = objectUrl;
    }

    function removeNewPhoto() {
        photoFile = null;
        portraitUrl = null;
    }

    async function onboardSeller() {
        isOnboarding = true;
        try {
            const res = await fetchWithAuth(`${$urlBase}/api/seller/onboard`, {
                method: 'POST'
            });

            if (res.status === 200) {
                const url = await res.json();
                window.location.href = url;
            } else {
                logger.error(`Failed to onboard seller: ${res.status}`);
                isOnboarding = false;
            }
        } catch (err) {
            logger.error(`Error onboarding seller: ${err}`);
            isOnboarding = false;
        }
    }

    async function saveChanges() {
        isSaving = true;
        saveError = null;
        saveSuccess = false;
        usernameError = null;

        if (username.includes(' ')) {
            usernameError = 'Username cannot contain spaces.';
            isSaving = false;
            return;
        }

        const form = new FormData();
        form.append('username', username);
        form.append('bio', bio);
        if (photoFile) {
            form.append('portraitFile', photoFile);
        }

        try {
            const res = await fetchWithAuth(`${$urlBase}/api/profile`, {
                method: "PUT",
                body: form
            });

            if (res.status === 200) {
                saveSuccess = true;
                usernameStore.set('');
                await loadProfile();
                username = get(usernameStore);
                bio = get(bioStore);
                portraitUrl = get(portraitUrlStore) || null;
                photoFile = null;
            } else if (res.status === 409) {
                usernameError = 'Username is already taken.';
            } else {
                logger.error(`Failed to save profile: ${res.status}`);
                saveError = 'Failed to save changes. Please try again.';
            }
        } catch (err) {
            logger.error(`Error saving profile: ${err}`);
            saveError = 'Technical issues, please try again later.';
        } finally {
            isSaving = false;
        }
    }
</script>

<main class={`transition-all duration-300 h-screen overflow-y-auto overflow-x-hidden w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>
    <TopHeader pageName="Profile" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center pb-32">
        <div class="outline-solid outline-zinc-600 rounded-3xl w-2/3 max-w-4xl flex flex-col items-center p-8">
            <h2 class="mb-8 text-3xl font-bold text-white">Edit Profile</h2>

            <!-- Email (read-only) -->
            <div class="w-full mb-6">
                <h3 class="text-xl font-semibold text-white mb-1">Email</h3>
                <p class="w-full px-3 py-2 rounded-lg bg-zinc-800 text-zinc-400 border border-zinc-700 select-none">
                    {$emailStore}
                </p>
            </div>

            <!-- Username -->
            <div class="w-full mb-6">
                <h3 class="text-xl font-semibold text-white mb-1">Username</h3>
                <input
                    class="w-full px-3 py-2 rounded-lg bg-zinc-700 text-white placeholder-zinc-400 border border-zinc-600 focus:border-indigo-500 focus:outline-hidden focus:ring-2 focus:ring-indigo-500/20"
                    type="text"
                    bind:value={username}
                    placeholder="Your username"
                    maxlength={30}
                />
                <p class="text-sm mt-1 text-zinc-400">{username.length}/30 characters</p>
                {#if usernameError}
                    <p class="text-sm mt-1 text-red-400">{usernameError}</p>
                {/if}
            </div>

            <!-- Bio -->
            <div class="w-full mb-6">
                <h3 class="text-xl font-semibold text-white mb-1">Bio</h3>
                <textarea
                    class="w-full px-3 py-2 rounded-lg bg-zinc-700 text-white placeholder-zinc-400 border border-zinc-600 focus:border-indigo-500 focus:outline-hidden focus:ring-2 focus:ring-indigo-500/20 resize-none"
                    rows={4}
                    bind:value={bio}
                    placeholder="Tell people a bit about yourself..."
                    maxlength={300}
                ></textarea>
                <p class="text-sm mt-1 text-zinc-400">{bio.length}/300 characters</p>
            </div>

            
            <!-- Profile Photo -->
            <div class="w-full mb-12 flex flex-col">
                <h3 class="text-xl font-semibold text-white mb-1">Profile Photo</h3>
                <div class="flex flex-row h-20 w-full items-center items-end gap-4">
                    <div class="h-20 w-20 shrink-0 rounded-xl bg-zinc-700 border border-zinc-600 overflow-hidden">
                        {#if portraitUrl}
                            <img src={portraitUrl} alt="Profile preview" class="h-full w-full object-cover" />
                        {/if}
                    </div>
                    <label for="photo" class="cursor-pointer h-10 ml-4 px-4 py-2 text-sm font-semibold bg-indigo-600 hover:bg-indigo-700 border border-zinc-600 rounded-lg text-white transition-colors flex items-center">
                        Change / Add Photo
                    </label>
                    <input id="photo" class="hidden" type="file" accept="image/*" onchange={handlePhotoChange} />
                </div>
            </div>

            <!-- Stripe Account Setup -->
            <div class="w-full mb-6">
                <h3 class="text-xl font-semibold text-white mb-1">Seller Account</h3>
                <p class="text-sm text-zinc-400 mb-3">Set up your Stripe account to start selling on Layerrs.</p>
                <button
                    onclick={onboardSeller}
                    disabled={isOnboarding}
                    class="px-6 py-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg text-white font-semibold text-sm transition-colors flex items-center gap-2"
                >
                    {#if isOnboarding}
                        <svg class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                        </svg>
                        Redirecting...
                    {:else}
                        Set Up Seller Account
                    {/if}
                </button>
            </div>


            <!-- Feedback -->
            {#if saveSuccess}
                <p class="mb-4 text-green-400 text-sm font-medium">Profile saved successfully.</p>
            {/if}
            {#if saveError}
                <p class="mb-4 text-red-400 text-sm font-medium">{saveError}</p>
            {/if}

            <!-- Save Button -->
            <button
                onclick={saveChanges}
                disabled={isSaving}
                class="px-8 py-4 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed rounded-full text-white font-semibold text-lg transition-colors"
            >{isSaving ? 'Saving...' : 'Save Changes'}</button>
        </div>
    </section>
</main>
