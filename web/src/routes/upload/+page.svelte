<script>
    import { goto } from "$app/navigation";


    import TopHeader from "../../components/TopHeader.svelte";
    import { jwt } from "../../stores/auth";
    import { isSidebarOpen } from "../../stores/player";

    let audioFiles = $state();
    let coverArtFiles = $state();
    let title = $state();
    let isUploaded = $state(false)

    // let data = $props();

    function removeAudioFile() {
        audioFiles = null
    }

    function removeCoverArtFile() {
        coverArtFiles = null
    }

    async function submitFile() {
        let audioFile = audioFiles[0];
        let coverArtFile = coverArtFiles[0];

        if (audioFile && coverArtFile) {
            console.log("audioFile and coverArtFile are valid")
            const form = new FormData();
            form.append('audioFile', audioFile)
            form.append('coverArtFile', coverArtFile)
            form.append('name', title)
            const res = await fetch("http://layerrs.com/api/track/", { 
                method: "POST", 
                headers: {
                    'Authorization': `Bearer ${$jwt}`
                },
                body: form
            });
            
            if (res.status == 201) {
                isUploaded = true
            }

        } else {
            console.log("audioFile and coverArtFile are not valid")
        }
    }

    function navigateHome() {
        goto('/')
    }

</script>

<main class={`transition-all duration-300 h-full w-full ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-gradient-to-b from-gray-800 to-gray-900`}>

    <TopHeader pageName="Upload" pageIcon="/upload.png"></TopHeader>

    <section class="w-full flex flex-row justify-center mt-12">

            <div class="outline outline-gray-600 rounded-3xl w-1/2 flex flex-col items-center">
                {#if !isUploaded}
                    <h2 class="my-8 text-4xl font-bold">Upload a Track</h2>
                    
                    {#if audioFiles}
                        <button onclick={removeAudioFile} class="my-4 w-1/2 bg-indigo-400 hover:bg-indigo-600 rounded-lg h-10 w-36">Remove Audio File</button>
                    {:else}
                        <input id="audio" class="hidden" type="file" accept="audio/*" bind:files={audioFiles} />
                        <label for="audio" class="my-4 w-1/2 bg-indigo-500 hover:bg-indigo-400 rounded-lg h-10 w-36 flex flex-row items-center justify-center">Upload Audio</label>
                    {/if}
                    
                    {#if coverArtFiles}
                        <button onclick={removeCoverArtFile} class="my-4 w-1/2 bg-indigo-600 hover:bg-indigo-400 rounded-lg h-10 w-36">Remove Cover Art</button>
                    {:else}
                        <input id="coverArt" class="hidden" type="file" accept="img/*" bind:files={coverArtFiles} />
                        <label for="coverArt" class="my-4 w-1/2 bg-indigo-500 hover:bg-indigo-400 rounded-lg h-10 w-36 flex flex-row items-center justify-center">Upload Cover Art</label>
                    {/if}
            
                    <input class="my-4 px-2 h-8 w-1/2 rounded-lg bg-indigo-200 hover:border-1 hover:border-black hover:bg-indigo-300" type="text" bind:value={title} placeholder="Track Name . . ." />
                    <button class="my-8 h-12 w-48 rounded-full bg-indigo-800 hover:bg-indigo-600" onclick={submitFile}>Upload</button>
                {:else}
                    <div class="w-full h-full flex flex-col items-center justify-center">
                        <h2 class="my-8 text-3xl font-bold">Track successfully uploaded!</h2>
                        <button class="bg-indigo-800 hover:bg-indigo-600 mb-12 h-12 w-36 rounded-full" onclick={navigateHome}>Return home</button>
                    </div>
                {/if}
            </div>
            
    </section>
</main>

