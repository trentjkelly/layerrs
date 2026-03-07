<script lang="ts">
    import { isLoggedIn } from "../stores/auth";
    import { goto } from "$app/navigation";
    import { handleBrowserLogout } from "../modules/lib/session";
    import { logger } from "../modules/lib/logger";

    let { username, email, portraitUrl } = $props();
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
        goto('/profile');
    }
</script>

<!-- Fixed circle in top right - only show when logged in -->
{#if $isLoggedIn}
<div class="fixed top-5 right-8 z-50">
    <div
        class="relative"
        onmouseenter={() => showDropdown = true}
        onmouseleave={() => showDropdown = false}
        role="button"
        tabindex="0"
    >
        {#if showDropdown}
            <!-- Expanded dropdown with circle integrated -->
            <div class="absolute top-0 right-0 w-56 bg-zinc-800 rounded-lg shadow-xl border border-zinc-700 z-50">
                <!-- Profile section with original circle -->
                <div class="px-4 py-4 border-b border-zinc-700">
                    <div class="flex items-center space-x-3">
                        <div class="w-12 h-12 bg-violet-700 rounded-xl flex-shrink-0 shadow-lg overflow-hidden">
                            {#if portraitUrl}
                                <img src={portraitUrl} alt="Profile" class="w-full h-full object-cover" />
                            {/if}
                        </div>
                        <div class="flex-1 min-w-0">
                            <p class="text-md font-medium text-white truncate">{username || 'Profile'}</p>
                            <p class="text-sm text-zinc-400 truncate">{email}</p>
                        </div>
                    </div>
                </div>

                <!-- Settings option -->
                <button
                    onclick={handleProfile}
                    class="w-full px-4 pt-3 pb-2 text-left text-zinc-300 hover:bg-zinc-700 flex items-center transition-colors duration-150"
                >
                    <svg class="w-4 h-4 mr-3 text-zinc-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                    </svg>
                    Profile
                </button>

                <!-- Logout option -->
                <button
                    onclick={handleLogout}
                    class="w-full px-4 py-2 text-left text-red-400 hover:bg-red-900/20 flex items-center transition-colors duration-150"
                >
                    <svg class="w-4 h-4 mr-3 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
                    </svg>
                    Logout
                </button>
            </div>
        {:else}
            <!-- Default circle when not expanded -->
            <div class="w-12 h-12 bg-violet-700 rounded-xl cursor-pointer shadow-lg hover:shadow-xl transition-all duration-200 hover:scale-105 overflow-hidden">
                {#if portraitUrl}
                    <img src={portraitUrl} alt="Profile" class="w-full h-full object-cover" />
                {/if}
            </div>
        {/if}
    </div>
</div>
{/if}
