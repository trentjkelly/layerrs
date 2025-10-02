<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
</script>

<main class={`transition-all duration-300 min-h-screen w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>
    <TopHeader pageName="Forgot Password" pageIcon=""></TopHeader>
    <div class="w-full h-full flex flex-col items-center justify-center">
        <h1 class="text-white text-4xl font-bold">Forgot Password</h1>
    </div>
</main>