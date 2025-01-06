<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import TrackCard from "../../components/TrackCard.svelte";
    import { isSidebarOpen } from "../../stores/player";

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }
</script>

<main class={`transition-all duration-300 h-full w-full p-8 ${$isSidebarOpen ? 'ml-64' : 'ml-0'}`}>

    <TopHeader pageName="Library"></TopHeader>

    <section>
        <TrackCard trackId={15}></TrackCard>  
    </section>

</main>
