<script>
    import TopHeader from "../components/TopHeader.svelte";
    import TrackCard from "../components/TrackCard.svelte";
    import { isSidebarOpen } from "../stores/player";

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }
</script>

<main class={`transition-all duration-300 h-full w-full p-8 ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-900 to-gray-800`}>

    <TopHeader pageName="Home" pageIcon="home.png"></TopHeader>

    <!-- Where the songs go -->
    <section>
        <!-- <div class="h-96 w-full flex flex-row justify-between">
            <TrackCard trackId={14}></TrackCard>  
            <TrackCard trackId={15}></TrackCard>  
            <TrackCard trackId={16}></TrackCard>  
            <TrackCard trackId={17}></TrackCard>  
        </div>
        <div class="h-96 w-full flex flex-row justify-between">
            <TrackCard trackId={14}></TrackCard>  
            <TrackCard trackId={15}></TrackCard>  
            <TrackCard trackId={16}></TrackCard>  
            <TrackCard trackId={17}></TrackCard>  
        </div> -->
    </section>
</main>


