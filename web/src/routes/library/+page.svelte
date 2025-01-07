<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import TrackCard from "../../components/TrackCard.svelte";
    import { isSidebarOpen } from "../../stores/player";

    function toggleSidebar() {
        $isSidebarOpen = !$isSidebarOpen;
    }
</script>

<main class={`transition-all duration-300 h-full w-full p-8 ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-800 to-gray-900`}>

    <TopHeader pageName="Your Library" pageIcon="vinyl.png"></TopHeader>

    <section class="w-full flex flex-wrap justify-around pb-24">
        <TrackCard trackId={14}></TrackCard>
        <TrackCard trackId={15}></TrackCard>
        <TrackCard trackId={16}></TrackCard>
        <TrackCard trackId={17}></TrackCard>
        <TrackCard trackId={14}></TrackCard>
        <TrackCard trackId={15}></TrackCard>
        <TrackCard trackId={16}></TrackCard>
        <TrackCard trackId={17}></TrackCard>
    </section>

</main>
