<script>
    import { isLoggedIn } from "../stores/auth";
    import { goto } from "$app/navigation";
    import { handleBrowserLogout } from "../modules/lib/session";
    import { logger } from "../modules/lib/logger";
    
    let showDropdown = false;

    async function handleLogout() {
        const success = await handleBrowserLogout();
        if (success) {
            goto('/login');
        } else {
            logger.error('Failed to logout');
        }
    }

    function handleSettings() {
        goto('/settings');
    }
</script>

<!-- Fixed circle in top right - only show when logged in -->
{#if $isLoggedIn}
<div class="fixed top-5 right-5 z-50">
    <div 
        class="relative"
        onmouseenter={() => showDropdown = true}
        onmouseleave={() => showDropdown = false}
        role="button"
        tabindex="0"
    >
        {#if showDropdown}
            <!-- Expanded dropdown with circle integrated -->
            <div class="absolute top-0 right-0 w-56 bg-gray-800 rounded-lg shadow-xl border border-gray-700 py-3 z-50">
                <!-- Profile section with original circle -->
                <div class="px-4 py-3 border-b border-gray-700">
                    <div class="flex items-center space-x-3">
                        <div class="w-12 h-12 bg-violet-700 rounded-full flex-shrink-0 shadow-lg"></div>
                        <div class="flex-1 min-w-0">
                            <p class="text-sm font-medium text-white truncate">User Profile</p>
                            <p class="text-xs text-gray-400 truncate">user@example.com</p>
                        </div>
                    </div>
                </div>
                
                <!-- Settings option -->
                <button 
                    onclick={handleSettings}
                    class="w-full px-4 py-2 text-left text-gray-300 hover:bg-gray-700 flex items-center transition-colors duration-150"
                >
                    <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    </svg>
                    Settings
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
            <div class="w-12 h-12 bg-violet-700 rounded-full cursor-pointer shadow-lg hover:shadow-xl transition-all duration-200 hover:scale-105">
            </div>
        {/if}
    </div>
</div>
{/if}
