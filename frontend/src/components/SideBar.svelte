<script>  
    import { isSidebarOpen } from "../stores/player";
    import { goto } from "$app/navigation";
    import { jwt } from "../stores/auth";

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

    function navigateSignIn() {
        goto('/login')
    }
    
</script>
  
<div class="flex">
    <div class={`fixed top-0 left-0 h-full bg-gradient-to-b from-gray-950 to-gray-900 text-white w-64 transform transition-transform duration-300 
        ${$isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div class="h-20 flex flex-row items-center justify-between">
            <button onclick={navigateHome} class="flex flex-row items-center ml-4 p-4">
                <img class="h-8 w-8" src="/vinyls.png" alt="Layerr Logo">
                <h2 class="text-2xl font-bold ml-2 tracking-wide">Layerrs</h2>
            </button>
            {#if $isSidebarOpen}
                <button class="mr-4 h-8 w-8 text-white rounded rounded-xl shadow hover:bg-gray-600 transform flex flex-row items-center justify-center" onclick={toggleSidebar}>
                    <img src="/sidebar-in.png" alt="Close sidebar"/>
                </button>
            {/if}
        </div>
        
        <ul class="space-y-2 px-4">
            <li class="hover:bg-slate-700 rounded">
                <button class="flex flex-row w-full h-full p-2" onclick={navigateHome}>
                    <img class="h-6 w-6" src="/home.png" alt="Home" />
                    <p class="pl-3 text-lg">Home</p>
                </button>
            </li>
            <li class="hover:bg-slate-700 rounded">
                <button class="flex flex-row w-full h-full p-2" onclick={navigateLibrary}>
                    <img class="h-6 w-6" src="/vinyl.png" alt="Library" />
                    <p class="pl-3 text-lg">Library</p>
                </button>
            </li>
            <li class="hover:bg-slate-700 rounded">
                <button class="flex flex-row w-full h-full p-2" onclick={navigateUpload}>
                    <img class="h-6 w-6" src="/upload.png" alt="Upload" />
                    <p class="pl-3 text-lg">Upload</p>
                </button>
            </li>
            { #if !$jwt }
                <li class="hover:bg-gray-500 rounded border text-indigo-200 border-indigo-200">
                    <button class="w-full h-full flex flex-row items-center justify-center p-2 " onclick={navigateSignIn}>Log In</button>
                </li>
            {/if}
        </ul>  
    </div>
</div>