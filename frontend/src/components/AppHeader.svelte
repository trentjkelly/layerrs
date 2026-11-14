<script lang="ts">
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import { emailStore, portraitUrlStore, usernameStore } from '../stores/profile';
    import UserMenu from './UserMenu.svelte';

    let { searchActive = $bindable(false) } = $props();

    let searchQuery = $state('');
    let searchInput: HTMLInputElement;
    let shortcutLabel = $state('Ctrl K');

    const navigation = [
        { label: 'Home', href: '/', icon: '/home.png' },
        { label: 'Library', href: '/library', icon: '/vinyl.png' },
        { label: 'Upload', href: '/upload', icon: '/upload.png' },
        { label: 'Pages', href: '/pages', icon: '/vinyl.png' }
    ];

    $effect(() => {
        searchQuery = $page.url.searchParams.get('search') ?? '';
    });

    function isActive(href: string) {
        if (href === '/') return $page.url.pathname === '/';
        return $page.url.pathname === href || $page.url.pathname.startsWith(`${href}/`);
    }

    function handleSearch(event: SubmitEvent) {
        event.preventDefault();
        const query = searchQuery.trim();
        if (!query) {
            searchInput.blur();
            return;
        }
        goto(`/?search=${encodeURIComponent(query)}`);
    }

    function focusSearch() {
        searchInput?.focus();
        searchInput?.select();
    }

    function toggleSearch() {
        if (document.activeElement === searchInput) {
            searchInput.blur();
            return;
        }

        focusSearch();
    }

    function handleGlobalShortcut(event: KeyboardEvent) {
        const isSearchShortcut =
            (event.metaKey || event.ctrlKey) &&
            !event.altKey &&
            !event.shiftKey &&
            event.key.toLowerCase() === 'k';

        if (!isSearchShortcut) return;

        event.preventDefault();
        event.stopPropagation();
        toggleSearch();
    }

    function handleSearchKeydown(event: KeyboardEvent) {
        if (event.key === 'Escape') {
            searchInput.blur();
        }
    }

    onMount(() => {
        const isApplePlatform = /Mac|iPhone|iPad|iPod/i.test(navigator.platform);
        shortcutLabel = isApplePlatform ? '⌘ K' : 'Ctrl K';

        window.addEventListener('keydown', handleGlobalShortcut, { capture: true });
        return () => window.removeEventListener('keydown', handleGlobalShortcut, { capture: true });
    });
</script>

<header class="relative z-40 shrink-0 border-b border-zinc-800 bg-zinc-900/95 backdrop-blur">
    <div class="mx-auto flex h-[4.75rem] w-full items-center gap-3 px-4 sm:gap-6 sm:px-6 lg:px-8">
        <a
            href="/"
            class={`shrink-0 font-humankind text-xl leading-none tracking-wider text-white transition-[color,filter] duration-150 hover:text-violet-300 sm:text-2xl ${searchActive ? 'blur-[1.5px]' : ''}`}
            aria-label="Layerrs home"
        >
            layerrs
        </a>

        <form class="mx-auto w-full max-w-2xl" role="search" onsubmit={handleSearch}>
            <label class="sr-only" for="global-search">Search tracks or artists</label>
            <div class="flex h-11 items-center border border-zinc-700 bg-zinc-950 transition-colors focus-within:border-violet-500">
                <svg class="ml-4 h-5 w-5 shrink-0 text-zinc-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true">
                    <circle cx="11" cy="11" r="7" stroke-width="2"></circle>
                    <path d="m20 20-4-4" stroke-width="2" stroke-linecap="round"></path>
                </svg>
                <input
                    id="global-search"
                    bind:this={searchInput}
                    bind:value={searchQuery}
                    type="search"
                    placeholder="Search tracks or artists"
                    aria-keyshortcuts="Meta+K Control+K"
                    onkeydown={handleSearchKeydown}
                    onfocus={() => (searchActive = true)}
                    onblur={() => (searchActive = false)}
                    class="search-input h-full min-w-0 flex-1 bg-transparent px-3 text-sm text-white outline-none placeholder:text-zinc-500 sm:text-base"
                />
                {#if searchQuery}
                    <span class="mr-3 hidden border border-violet-500/70 px-1.5 py-0.5 text-[0.65rem] text-violet-400 sm:block" aria-hidden="true">
                        ENTER
                    </span>
                {:else}
                    <button
                        type="button"
                        onclick={focusSearch}
                        aria-label="Focus search"
                        aria-keyshortcuts="Meta+K Control+K"
                        class="mr-3 hidden border border-zinc-700 px-1.5 py-0.5 text-[0.65rem] text-zinc-500 transition-colors hover:border-zinc-500 hover:text-zinc-300 sm:block"
                    >
                        {shortcutLabel}
                    </button>
                {/if}
            </div>
        </form>

        <div class={`shrink-0 transition-[filter] duration-150 ${searchActive ? 'blur-[1.5px]' : ''}`}>
            <UserMenu username={$usernameStore} email={$emailStore} portraitUrl={$portraitUrlStore} />
        </div>
    </div>

    <nav
        aria-label="Primary navigation"
        class={`mx-auto w-full px-2 transition-[filter] duration-150 sm:px-6 lg:px-8 ${searchActive ? 'blur-[1.5px]' : ''}`}
    >
        <ul class="grid w-full grid-cols-4 items-end gap-1 sm:mx-auto sm:max-w-2xl sm:gap-3">
            {#each navigation as item}
                <li class="min-w-0">
                    <a
                        href={item.href}
                        aria-current={isActive(item.href) ? 'page' : undefined}
                        class={`group relative flex h-12 min-w-0 items-center justify-center gap-2 px-1 text-xs font-medium transition-colors sm:px-4 sm:text-sm ${isActive(item.href) ? 'text-white' : 'text-zinc-400 hover:text-white'}`}
                    >
                        <img src={item.icon} alt="" class={`hidden h-4 w-4 transition-opacity sm:block ${isActive(item.href) ? 'opacity-100' : 'opacity-60 group-hover:opacity-100'}`} />
                        <span>{item.label}</span>
                        <span class={`absolute inset-x-3 bottom-0 h-0.5 bg-violet-500 transition-opacity ${isActive(item.href) ? 'opacity-100' : 'opacity-0'}`}></span>
                    </a>
                </li>
            {/each}
        </ul>
    </nav>
</header>

<style>
    .search-input::-webkit-search-cancel-button {
        -webkit-appearance: none;
        appearance: none;
    }
</style>
