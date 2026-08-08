<script lang="ts">
    import { onMount } from 'svelte';
    import { isLoggedIn } from '../stores/auth';
    import { artistIdStore, loadProfile } from '../stores/profile';
    import { addTrackToPage, getArtistPages } from '../modules/requests/page-requests';
    import type { Page } from '../models/types';

    let { trackId, urlBase }: { trackId: number; urlBase: string } = $props();

    let showMenu = $state(false);
    let showModal = $state(false);
    let pages = $state<Page[]>([]);
    let isLoadingPages = $state(false);
    let addingPageId = $state<number | null>(null);
    let errorMessage = $state('');
    let successMessage = $state('');

    let menuBtnEl = $state<HTMLElement | null>(null);
    let menuPopupEl = $state<HTMLElement | null>(null);

    onMount(() => {
        loadProfile();
    });

    function toggleMenu(event: MouseEvent) {
        event.stopPropagation();
        showMenu = !showMenu;
    }

    function closeMenu() {
        showMenu = false;
    }

    function handleClickOutside(event: MouseEvent) {
        const target = event.target as Node;
        if (
            showMenu &&
            menuPopupEl &&
            !menuPopupEl.contains(target) &&
            menuBtnEl &&
            !menuBtnEl.contains(target)
        ) {
            closeMenu();
        }
    }

    $effect(() => {
        if (showMenu) {
            document.addEventListener('click', handleClickOutside);
            return () => document.removeEventListener('click', handleClickOutside);
        }
    });

    async function openPagePicker() {
        closeMenu();
        showModal = true;
        errorMessage = '';
        successMessage = '';
        await loadPages();
    }

    async function loadPages() {
        let artistId = $artistIdStore;
        if (!artistId) {
            await loadProfile();
            artistId = $artistIdStore;
        }
        if (!artistId) {
            isLoadingPages = false;
            return;
        }

        isLoadingPages = true;
        pages = await getArtistPages(urlBase, String(artistId));
        isLoadingPages = false;
    }

    async function handleSelectPage(page: Page) {
        addingPageId = page.id;
        errorMessage = '';
        successMessage = '';
        const result = await addTrackToPage(urlBase, String(page.id), trackId);
        addingPageId = null;

        if (result) {
            successMessage = `Added to ${page.name}`;
            setTimeout(() => {
                showModal = false;
                successMessage = '';
            }, 1200);
        } else {
            errorMessage = 'Could not add track to page.';
        }
    }

    function closeModal() {
        showModal = false;
        errorMessage = '';
        successMessage = '';
    }
</script>

{#if $isLoggedIn}
    <div class="relative">
        <button
            bind:this={menuBtnEl}
            class="mr-1 p-1 hover:bg-olive-600 hover:cursor-pointer flex flex-row items-center"
            onclick={toggleMenu}
            aria-label="More"
        >
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                <circle cx="12" cy="6" r="1.5" />
                <circle cx="12" cy="12" r="1.5" />
                <circle cx="12" cy="18" r="1.5" />
            </svg>
        </button>

        {#if showMenu}
            <div
                bind:this={menuPopupEl}
                class="absolute top-full right-0 mt-1 w-56 bg-olive-500 shadow-xl border border-white z-50"
            >
                <button
                    class="w-full px-4 py-2 text-left text-white hover:bg-olive-600 flex items-center transition-colors duration-150 cursor-pointer"
                    onclick={openPagePicker}
                >
                    <svg class="w-4 h-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                        />
                    </svg>
                    Add to page
                </button>
            </div>
        {/if}
    </div>

    {#if showModal}
        <div
            class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
            role="button"
            tabindex="-1"
            aria-label="Close"
            onclick={(e) => e.target === e.currentTarget && closeModal()}
            onkeydown={(e) => e.key === 'Escape' && closeModal()}
        >
            <div class="w-full max-w-md bg-olive-500 border border-white shadow-xl">
                <div class="flex items-center justify-between px-4 py-3 border-b border-white/30">
                    <h3 class="text-lg font-medium text-white">Add track to a page</h3>
                    <button
                        onclick={closeModal}
                        class="text-white hover:text-olive-200 cursor-pointer"
                        aria-label="Close"
                    >
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>

                <div class="px-4 py-4 max-h-[60vh] overflow-y-auto">
                    {#if isLoadingPages}
                        <p class="text-white text-center py-4">Loading your pages…</p>
                    {:else if pages.length === 0}
                        <p class="text-white text-center py-4">You don't have any pages yet.</p>
                        <a
                            href="/pages/new"
                            class="block text-center text-white underline hover:text-olive-200"
                        >
                            Create a page
                        </a>
                    {:else}
                        <div class="flex flex-col gap-2">
                            {#each pages as page (page.id)}
                                <button
                                    onclick={() => handleSelectPage(page)}
                                    disabled={addingPageId !== null}
                                    class="w-full text-left px-3 py-2 text-white hover:bg-olive-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-150 cursor-pointer"
                                >
                                    {#if addingPageId === page.id}
                                        Adding…
                                    {:else}
                                        {page.name}
                                    {/if}
                                </button>
                            {/each}
                        </div>
                    {/if}

                    {#if errorMessage}
                        <p class="mt-3 text-red-300 text-sm">{errorMessage}</p>
                    {/if}
                    {#if successMessage}
                        <p class="mt-3 text-green-300 text-sm">{successMessage}</p>
                    {/if}
                </div>
            </div>
        </div>
    {/if}
{/if}
