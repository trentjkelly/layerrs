<script>

    let files = $state();

    async function submitFile() {
        let selectedFile = files[0];

        if (selectedFile) {
            console.log("There is a selected file")
            const form = new FormData();
            form.append('file', selectedFile)

            const res = await fetch("http://localhost:8080/track", {method: "PUT", body: form});
            console.log(res)
        } else {
            console.log("No selected file")
        }
    }
    
</script>

<div class="h-screen w-screen flex flex-row">

    <input type="file" accept="audio/*" bind:files />
    <button onclick={submitFile}>Submit</button>
   
</div>
