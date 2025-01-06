<script>
    import { isSidebarOpen } from "../stores/player";

    let { pageName, pageIcon } = $props();

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }

</script>

<section class="h-24 w-full flex flex-row">

    {#if !$isSidebarOpen}
        <button onclick={toggleSidebar} class="h-8 w-8 hover:h-9 hover:w-9 text-white rounded rounded-xl shadow hover:bg-gray-800">
            <img src="right-arrow.png" alt="Open sidebar">
        </button>    
    {/if}
    <div class="text-3xl font-bold ml-8 flex flex-row">
        {#if pageIcon !== ""}
            <img class="h-8 w-8 mr-4" src={pageIcon} alt="Home Icon">
        {/if}
        <h1>{ pageName }</h1>
    </div>

</section>
