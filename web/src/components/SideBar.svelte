<script>  
    import { isSidebarOpen } from "../stores/player";
    import { goto } from "$app/navigation";

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
        goto('/signin')
    }

</script>
  
<div class="flex">
    <div
      class={`fixed top-0 left-0 h-full bg-slate-700 text-white w-64 transform transition-transform duration-300 ${
        $isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div class="flex flex-row items-center justify-between">
            <h2 class="text-xl font-bold p-4">Layerr</h2>
            {#if $isSidebarOpen}
                <button class="mr-4 h-8 w-12 bg-blue-500 text-white rounded shadow hover:bg-blue-600 transform" onclick={toggleSidebar}>
                    Close
                </button>
            {/if}
        </div>
        
        <ul class="space-y-2 px-4">
            <li class="hover:bg-slate-500 p-2 rounded" onclick={navigateHome}>Home</li>
            <li class="hover:bg-gray-500 p-2 rounded" onclick={navigateLibrary}>Library</li>
            <li class="hover:bg-gray-500 p-2 rounded" onclick={navigateUpload}>Upload</li>
            <li class="hover:bg-gray-500 p-2 rounded border text-indigo-200 border-indigo-200" onclick={navigateSignIn}>Sign-in</li>
        </ul>
  
    </div>
</div>
  