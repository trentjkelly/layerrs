<script>
    import { isSidebarOpen } from "../stores/player";

    let { pageName } = $props();

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }
</script>

<section class="h-24 w-full flex flex-row">

    {#if !$isSidebarOpen}
        <button onclick={toggleSidebar} class="h-8 w-12 bg-blue-500 text-white rounded shadow hover:bg-blue-600">Open</button>    
    {/if}
    <h1 class="text-3xl font-bold ml-8">{ pageName }</h1>

</section>