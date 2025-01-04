<script>

    let audioFiles = $state();
    let coverArtFiles = $state();
    let title = $state();
    let artistId = $state();

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
            form.append('artistId', artistId)
            const res = await fetch("http://localhost:8080/api/track/", {method: "POST", body: form});
            console.log(res)
        } else {
            console.log("audioFile and coverArtFile are not valid")
        }
    }

</script>

<div class="h-screen w-screen flex flex-row justify-center items-center">
    <div class="border-4 rounded-3xl w-1/2 flex flex-col items-center justify-center">
        <h2 class="my-8 text-5xl text-indigo-700">Upload a Track</h2>
        
        {#if audioFiles}
            <button onclick={removeAudioFile} class="my-4 w-1/2 text-white bg-indigo-600 hover:bg-indigo-400 rounded-lg h-10 w-36">Remove Audio File</button>
        {:else}
            <input id="audio" class="hidden" type="file" accept="audio/*" bind:files={audioFiles} />
            <label for="audio" class="my-4 w-1/2 bg-indigo-200 hover:bg-indigo-400 rounded-lg h-10 w-36 flex flex-row items-center justify-center">Upload Audio</label>
        {/if}
        
        {#if coverArtFiles}
            <button onclick={removeCoverArtFile} class="my-4 w-1/2 text-white bg-indigo-600 hover:bg-indigo-400 rounded-lg h-10 w-36">Remove Cover Art</button>
        {:else}
            <input id="coverArt" class="hidden" type="file" accept="img/*" bind:files={coverArtFiles} />
            <label for="coverArt" class="my-4 w-1/2 bg-indigo-200 hover:bg-indigo-400 rounded-lg h-10 w-36 flex flex-row items-center justify-center">Upload Cover Art</label>
        {/if}

        <input class="my-4 px-2 h-8 w-1/2 rounded-lg bg-indigo-200 hover:border-1 hover:border-black hover:bg-indigo-300" type="text" bind:value={title} placeholder="Track Name . . ." />
        <input class="my-4 px-2 h-8 w-1/2 rounded-lg bg-indigo-200 hover:border-1 hover:border-black hover:bg-indigo-300" type="text" bind:value={artistId} placeholder="Enter artistId . . ." />
        <button class="my-8 h-12 w-48 rounded-full bg-indigo-500 hover:bg-indigo-700 text-white" onclick={submitFile}>Upload</button>
    </div>
</div>
