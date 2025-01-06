<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
</script>

<main class={`transition-all duration-300 h-full w-full p-8 ${$isSidebarOpen ? 'ml-64' : 'ml-0'}`}>

    <TopHeader pageName="Log in" pageIcon=""></TopHeader>

    <section>

    </section>
</main>