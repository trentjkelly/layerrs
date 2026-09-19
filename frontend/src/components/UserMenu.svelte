<script lang="ts">
    import { isLoggedIn } from "../stores/auth";
    import { goto } from "$app/navigation";
    import { handleBrowserLogout } from "../modules/lib/session";
    import { logger } from "../modules/lib/logger";

    let { username = '', email = '', portraitUrl = '' } = $props();
    let showDropdown = $state(false);

    async function handleLogout() {
        const success = await handleBrowserLogout();
        if (success) {
            goto('/login');
        } else {
            logger.error('Failed to logout');
        }
    }

    function handleProfile() {
        showDropdown = false;
        goto('/profile');
    }

    function handleLogin() {
        goto('/login');
    }

    function handleFocusOut(event: FocusEvent) {
        const nextTarget = event.relatedTarget as Node | null;
        if (!nextTarget || !(event.currentTarget as HTMLElement).contains(nextTarget)) {
            showDropdown = false;
        }
    }
</script>

{#if $isLoggedIn}
<div class="relative shrink-0" onfocusout={handleFocusOut}>
    <div
        class="relative"
    >
        <button
            type="button"
            onclick={() => showDropdown = !showDropdown}
            aria-label="Open profile menu"
            aria-haspopup="menu"
            aria-expanded={showDropdown}
            class="flex h-11 w-11 items-center justify-center overflow-hidden rounded-full border border-zinc-700 bg-zinc-700 text-sm font-bold text-white transition-colors hover:border-zinc-500 hover:bg-zinc-600 focus:outline-none focus:ring-2 focus:ring-violet-500 focus:ring-offset-2 focus:ring-offset-zinc-900"
        >
            {#if portraitUrl}
                <img src={portraitUrl} alt="" class="h-full w-full object-cover" />
            {:else}
                <span>{(username || email || 'P').slice(0, 1).toUpperCase()}</span>
            {/if}
        </button>

        {#if showDropdown}
            <div class="absolute right-0 top-[calc(100%+0.75rem)] z-50 w-64 border border-zinc-700 bg-zinc-800 shadow-2xl" role="menu">
                <div class="border-b border-zinc-700 px-4 py-4">
                    <p class="truncate text-sm font-semibold text-white">{username || 'Your profile'}</p>
                    <p class="mt-0.5 truncate text-sm text-zinc-400">{email}</p>
                </div>

                <button
                    onclick={handleProfile}
                    role="menuitem"
                    class="flex w-full cursor-pointer items-center px-4 py-3 text-left text-sm text-white transition-colors hover:bg-zinc-700"
                >
                    <svg class="mr-3 h-4 w-4 text-zinc-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                    </svg>
                    Profile
                </button>

                <button
                    onclick={handleLogout}
                    role="menuitem"
                    class="flex w-full cursor-pointer items-center px-4 py-3 text-left text-sm text-white transition-colors hover:bg-zinc-700"
                >
                    <svg class="mr-3 h-4 w-4 text-zinc-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
                    </svg>
                    Logout
                </button>
            </div>
        {/if}
    </div>
</div>
{:else}
    <button class="h-11 shrink-0 border border-zinc-600 px-3 text-sm font-semibold text-white transition-colors hover:border-zinc-400 hover:bg-zinc-800 sm:px-4" onclick={handleLogin}>Log in</button>
{/if}
