<script lang="ts">
	import '../app.css';
    import AudioPlayer from '../components/AudioPlayer.svelte';
    import SideBar from '../components/SideBar.svelte';
	let { children } = $props();
</script>

<div class="h-screen w-screen flex flex-row bg-gray-900 text-white">

	<SideBar></SideBar>

	{@render children()}

	<AudioPlayer></AudioPlayer>
    
</div>
