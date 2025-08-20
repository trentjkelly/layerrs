<script>
    import { onMount } from "svelte";
    import TopHeader from "../components/TopHeader.svelte";
    import NewTrackCard from "../components/NewTrackCard.svelte";
    import UserMenu from "../components/UserMenu.svelte";
    import { isSidebarOpen } from "../stores/player";
    import { handleEnvironment, urlBase } from "../stores/environment";

    // Each of the songs to be loaded in
    let artistId = 20 // Static for now

    /**
     * @type {any[]}
     */
    let trackIds = [];

    async function fetchData() {
        const response = await fetch(`${$urlBase}/api/recommendations/home`)
        const data = await response.json();
        trackIds = Object.keys(data).map(key => data[key])
    }

    onMount(async () => {
        await handleEnvironment()
        await fetchData()
    })

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }

</script>

<main class={`transition-all duration-300 h-auto w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-br from-gray-900 to-gray-800`}>

    <TopHeader pageName="Home" pageIcon="home.png"></TopHeader>

    <!-- User Menu -->
    <UserMenu />

    <!-- Where the songs go -->
    <section class="w-full flex flex-wrap justify-around pb-24">
        <!-- {#each trackIds as id}
            <NewTrackCard trackId={id}></NewTrackCard>
        {/each} -->
        <NewTrackCard trackId={artistId}></NewTrackCard>
        <NewTrackCard trackId={artistId}></NewTrackCard>
    </section>
</main>
