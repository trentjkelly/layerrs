<script lang="ts">
    import { onMount } from 'svelte';
	import '../app.css';
    import AudioPlayer from '../components/AudioPlayer.svelte';
    import SideBar from '../components/SideBar.svelte';
    import { audio } from '../stores/player';
	let { children } = $props();

	// Initialize audio component across the entire session
	onMount(() => {
		audio.set(new Audio())
	})

</script>

<div class="h-screen w-screen flex flex-row bg-gray-900 text-white">

	<SideBar></SideBar>

	{@render children()}

	<AudioPlayer></AudioPlayer>
    
</div>
