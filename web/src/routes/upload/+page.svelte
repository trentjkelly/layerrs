<script>

    let audioFiles = $state();
    let coverArtFiles = $state();
    let title = $state();
    let artistId = $state();

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
        <h2 class="my-4 text-3xl text-indigo-700">Upload a Track</h2>
        <input class="my-4 w-1/2 file:bg-indigo-200 hover:file:bg-indigo-400 file:border-0 file:rounded-2xl file:h-10 file:w-36 file:mr-4" type="file" accept="audio/*" bind:files={audioFiles} />
        <input class="my-4 w-1/2 file:bg-indigo-200 hover:file:bg-indigo-400 file:border-0 file:rounded-2xl file:h-10 file:w-36 file:mr-4   " type="file" accept="img/*" bind:files={coverArtFiles} />
        <input class="my-2 px-2 h-8 w-1/2 rounded-lg bg-indigo-200 hover:border-1 hover:border-black hover:bg-indigo-300" type="text" bind:value={title} placeholder="Track Name . . ." />
        <input class="my-2 px-2 h-8 w-1/2 rounded-lg bg-indigo-200 hover:border-1 hover:border-black hover:bg-indigo-300" type="text" bind:value={artistId} placeholder="Enter artistId . . ." />
        <button class="my-8 h-12 w-48 rounded-full bg-indigo-500 hover:bg-indigo-700 text-white" onclick={submitFile}>Upload</button>
    </div>
</div>
