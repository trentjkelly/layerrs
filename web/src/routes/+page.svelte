<script>
    // let trackId = "offsneakbeats-thank-god-ii.mp3";
    // let audioSrc = `http://localhost:8080/audio/${trackId}`;

    // let coverId = "kanye.png"
    // let coverSrc = `http://localhost:8080/covers/${coverId}`

    let files = $state();

    function submitFile() {
        let selectedFile = files[0];

        const form = new FormData();
        form.append('file', selectedFile)

        
    }

    // function handleFileChange({ target }) {
    //     const file = target.files[0]

    //     if (file) {
    //         selectedFile = file
    //     } else {
    //         console.log("No file found")
    //     }
    // }
    
</script>

<div class="h-screen w-screen flex flex-row">

    <!-- <nav class="h-screen w-80 text-xl bg-black text-white font-bold flex flex-col">
        <button>Layerr</button>
        <button>Profile</button>
        <button>Search</button>
        <button>Trending</button>
    </nav>

    <main class="w-full h-screen flex flex-row justify-center items-center">
        <img class="h-64 w-64" src={coverSrc} alt="Graduation album cover">
        <audio src={audioSrc} controls>
            Your browser does not support the audio element
        </audio>
    </main> -->

    <!-- <input type="file" accept="audio/*" bind:files />
    <button onclick={submitFile}>Submit</button> -->

    
</div>
