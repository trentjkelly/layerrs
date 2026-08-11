<script lang="ts">
	import { isSidebarOpen } from '../../../stores/player';
	import TopHeader from '../../../components/TopHeader.svelte';
 </script>

<main class={`transition-all duration-300 h-screen w-full overflow-y-auto ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-950`}>
	<TopHeader pageName="Track" pageIcon="" />
</main>
