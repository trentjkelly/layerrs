<script>
    import { isSidebarOpen } from "../../../stores/player";
    import TopHeader from "../../../components/TopHeader.svelte";
    import TrackCard from "../../../components/TrackCard.svelte";
    import TrackRow from "../../../components/TrackRow.svelte";

    let artistName = "Trent Kelly"
    let artistUsername = "trentjkelly"
    let artistBio = "Insert some corny ass drake quote as your bio ig"

    async function getArtistData() {
        // const backendURL = import.meta.env.VITE_BACKEND_URL;
        await fetch(`https://layerrs.com/api/artist/`)
    }

</script>

<main class={`transition-all duration-300 h-full w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-800 to-gray-900`}>

    <TopHeader pageName="" pageIcon=""></TopHeader>

    <section class="h-full w-full flex flex-col items-center">
        <div class="w-2/3 h-auto flex flex-row">
            <div class="w-48 h-48 bg-slate-600">
                <img src="" alt="">
            </div>
    
            <div class="ml-8">
                <h1 class="text-4xl font-bold">{artistName}</h1>
                <h2 class="text-xl">@{artistUsername}</h2>
                <p class="mt-4">{artistBio}</p>
            </div>
        </div>
        <div class="h-auto w-3/4 flex flex-col mt-8">
            <h2 class="text-2xl font-bold mb-4">Top Tracks</h2>
            <div class="flex flex-wrap w-full justify-around">
                <TrackRow />
            </div>
        </div>
    </section>

</main>