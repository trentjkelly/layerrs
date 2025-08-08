<script>
    import { isSidebarOpen } from "../stores/player";

    let { pageName, pageIcon } = $props();

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }

</script>

<section class="h-20 w-full flex flex-row items-center pl-8">

    {#if !$isSidebarOpen}
        <button onclick={toggleSidebar} class="h-8 w-8 text-white rounded rounded-xl shadow hover:bg-gray-600">
            <img src="/sidebar-out.png" alt="Open sidebar">
        </button>    
    {/if}
    <div class="text-2xl font-bold ml-8 flex flex-row items-center">
        {#if pageIcon !== ""}
            <img class="h-6 w-6 mr-2" src={pageIcon} alt="Home Icon">
        {/if}
        <h1>{ pageName }</h1>
    </div>

</section>
