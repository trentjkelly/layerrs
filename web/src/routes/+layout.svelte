<script lang="ts">
	import '../app.css';
    import AudioPlayer from '../components/AudioPlayer.svelte';
	let { children } = $props();
</script>

<div class="h-screen w-screen flex flex-row">

    <nav class="h-screen w-80 text-xl bg-black text-white font-bold flex flex-col">
        <button>Layerr</button>
        <button>Profile</button>
        <button>Search</button>
        <button>Trending</button>
        <button>Upload</button>
    </nav>

	{@render children()}

	<AudioPlayer></AudioPlayer>
    
</div>
