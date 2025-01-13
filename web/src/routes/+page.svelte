<script>
    import { onMount } from "svelte";
    import TopHeader from "../components/TopHeader.svelte";
    import TrackCard from "../components/TrackCard.svelte";
    import { isSidebarOpen } from "../stores/player";
    import { jwt } from "../stores/auth";

    // Set session jwt to the cookie jwt if it exists
    export let data 
    if (data.token) {
        jwt.set(data.token)
    }

    // Each of the songs to be loaded in
    let artistId = 1 // Static for now

    /**
     * @type {any[]}
     */
    let trackIds = [];


    async function fetchData() {
        const response = await fetch(`http://localhost:8080/api/recommendations/home`)
        const data = await response.json();
        trackIds = Object.keys(data).map(key => data[key])
    }

    onMount(() => {
        fetchData()
    })

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }

</script>

<main class={`transition-all duration-300 h-auto w-full p-8 ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-800 to-gray-900`}>

    <TopHeader pageName="Home" pageIcon="home.png"></TopHeader>

    <!-- Where the songs go -->
    <section class="w-full flex flex-wrap justify-around pb-24">
        {#each trackIds as id}
            <TrackCard trackId={id}></TrackCard>
        {/each}
    </section>
</main>
