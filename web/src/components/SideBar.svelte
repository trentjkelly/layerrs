<script>  
    import { isSidebarOpen } from "../stores/player";

    function toggleSidebar() {
      $isSidebarOpen = !$isSidebarOpen;
    }
</script>
  
<div class="flex">
    <div
      class={`fixed top-0 left-0 h-full bg-slate-700 text-white w-64 transform transition-transform duration-300 ${
        $isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div class="flex flex-row items-center justify-between">
            <h2 class="text-xl font-bold p-4">Layerr</h2>
            {#if $isSidebarOpen}
                <button class="mr-4 h-8 w-12 bg-blue-500 text-white rounded shadow hover:bg-blue-600 transform" on:click={toggleSidebar}>
                    Close
                </button>
            {/if}
        </div>
        
        <ul class="space-y-2 px-4">
            <li class="hover:bg-slate-500 p-2 rounded">Home</li>
            <li class="hover:bg-gray-500 p-2 rounded">Library</li>
            <li class="hover:bg-gray-500 p-2 rounded">Upload</li>
            <li class="hover:bg-gray-500 p-2 rounded border text-indigo-200 border-indigo-200">Sign-in</li>
        </ul>
  
    </div>
</div>
  