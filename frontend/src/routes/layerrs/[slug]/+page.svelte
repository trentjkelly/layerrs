<script>
    import TopHeader from "../../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../../stores/player";
</script>

<main class={`transition-all duration-300 h-full w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>

    <TopHeader pageName="" pageIcon=""></TopHeader>
    <p class="text-white">Layerrs</p>
</main>