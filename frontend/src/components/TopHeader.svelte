<script>
    import { isSidebarOpen } from "../stores/player";

    let { pageName, pageIcon } = $props();

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }

</script>

<section class="h-20 w-full flex flex-row items-center pl-8">

    {#if !$isSidebarOpen}
        <button onclick={toggleSidebar} class="h-8 w-8 text-white flex flex-row items-center justify-center hover:bg-zinc-500 hover:cursor-pointer">
            <img src="/sidebar-out.png" alt="Open sidebar" class="h-4 w-4">
        </button>    
    {/if}
    <div class="text-2xl font-bold ml-8 flex flex-row items-center">
        {#if pageIcon !== ""}
            <img class="h-6 w-6 mr-2" src={pageIcon} alt="Home Icon">
        {/if}
        <h1>{ pageName }</h1>
    </div>

</section>
