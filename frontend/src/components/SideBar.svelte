<script>  
    import { isSidebarOpen } from "../stores/player";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";

    function toggleSidebar() {
      $isSidebarOpen = !$isSidebarOpen;
    }

    function navigateHome() {
        goto('/')
    }

    function navigateLibrary() {
        goto('/library')
    }

    function navigateUpload() {
        goto('/upload')
    }

    function navigatePages() {
        goto('/pages')
    }

</script>
  
<div class="flex">
    <div class={`fixed top-0 left-0 h-full bg-zinc-900 text-white w-64 transform transition-transform duration-300
        ${$isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div class="h-20 flex flex-row items-center justify-between">
            <button onclick={navigateHome} class="flex flex-row ml-4 p-4 hover:cursor-pointer">
                <!-- <img class="h-8 w-8" src="/vinyls.png" alt="Layerr Logo"> -->
                <h2 class="text-2xl font-humankind font-space tracking-wider leading-none">layerrs</h2>
            </button>
            {#if $isSidebarOpen}
                <button class="mr-4 h-8 w-8 hover:cursor-pointer text-white hover:bg-zinc-800 transform flex flex-row items-center justify-center" onclick={toggleSidebar}>
                    <img src="/sidebar-in.png" alt="Close sidebar" class="h-4 w-4" />
                </button>
            {/if}
        </div>
        
        <ul class="space-y-1 px-4">
            <li class="hover:bg-zinc-800" class:bg-zinc-800={$page.url.pathname === '/'}>
                <button class="flex flex-row hover:cursor-pointer w-full h-full p-2" onclick={navigateHome}>
                    <img class="h-6 w-6" src="/home.png" alt="Home" />
                    <p class="pl-3 text-lg">Home</p>
                </button>
            </li>
            <li class="hover:bg-zinc-800" class:bg-zinc-800={$page.url.pathname === '/library'}>
                <button class="flex flex-row hover:cursor-pointer w-full h-full p-2" onclick={navigateLibrary}>
                    <img class="h-6 w-6" src="/vinyl.png" alt="Library" />
                    <p class="pl-3 text-lg">Library</p>
                </button>
            </li>
            <li class="hover:bg-zinc-800" class:bg-zinc-800={$page.url.pathname === '/upload'}>
                <button class="flex flex-row hover:cursor-pointer w-full h-full p-2" onclick={navigateUpload}>
                    <img class="h-6 w-6" src="/upload.png" alt="Upload" />
                    <p class="pl-3 text-lg">Upload</p>
                </button>
            </li>
            <li class="hover:bg-zinc-800" class:bg-zinc-800={$page.url.pathname === '/pages' || $page.url.pathname.startsWith('/pages/')}>
                <button class="flex flex-row hover:cursor-pointer w-full h-full p-2" onclick={navigatePages}>
                    <img class="h-6 w-6" src="/vinyl.png" alt="Pages" />
                    <p class="pl-3 text-lg">Pages</p>
                </button>
            </li>

        </ul>  
    </div>
</div>
