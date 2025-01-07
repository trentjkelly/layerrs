<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import { isSidebarOpen } from "../../stores/player";
</script>

<main class={`transition-all duration-300 h-full w-full p-8 ${$isSidebarOpen ? 'ml-64' : 'ml-0'}`}>

    <TopHeader pageName="Sign in" pageIcon=""></TopHeader>

    <section class="w-full flex flex-row justify-center">

        <div class="border border-gray-600 bg-gray-800 rounded-3xl w-1/2 flex flex-col items-center">
            <h2 class="my-8 text-4xl">Sign In</h2>
        </div>
        
    </section>
</main>