<script>
    import TopHeader from "../../components/TopHeader.svelte";
    import TrackCard from "../../components/TrackCard.svelte";
    import { isSidebarOpen } from "../../stores/player";
    import { onMount } from "svelte";
    import { jwt } from "../../stores/auth";
    import { isLoggedIn } from "../../stores/auth";
    import LogInPopup from "../../components/LogInPopup.svelte";

    /**
     * @type {any[]}
     */
    let trackIds = [];

    async function fetchData() {
        console.log("jwt" + $jwt)

        const response = await fetch(`http://localhost:8080/api/recommendations/library`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${$jwt}`
            }
        })

        if (response.status == 401) {
            console.log("unauthorized request for library")
            return
        }

        const data = await response.json();
        console.log(data)
        trackIds = Object.keys(data).map(key => data[key])
    }

    onMount(async () => {
        console.log()
        await fetchData()
    })

</script>

<main class={`transition-all duration-300 h-full w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-800 to-gray-900`}>

    <TopHeader pageName="Your Library" pageIcon="/vinyl.png"></TopHeader>

    
    <section class="w-full flex flex-wrap justify-around pb-24">
        {#if $isLoggedIn}

            {#each trackIds as id}
                {#if id}
                    <TrackCard trackId={id}></TrackCard>
                {/if}
            {/each}

        {:else}
            <LogInPopup />
        {/if}
    </section>
</main>
