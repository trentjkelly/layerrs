<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { handleBrowserLogin } from '../../../modules/lib/session';
	import { logger } from '../../../modules/lib/logger';

	let status = $state<'loading' | 'success' | 'error'>('loading');
	let errorMessage = $state('');

	onMount(() => {
		const hash = typeof window !== 'undefined' ? window.location.hash.slice(1) : '';
		const params = new URLSearchParams(hash);
		const jwt = params.get('jwt');
		const refresh = params.get('refresh');

		if (!jwt || !refresh) {
			logger.error('Callback missing jwt or refresh in URL hash');
			status = 'error';
			errorMessage = 'Invalid sign-in link. Please request a new one from the login page.';
			return;
		}

		handleBrowserLogin(jwt, refresh).then((success) => {
			if (success) {
				status = 'success';
				goto('/');
			} else {
				status = 'error';
				errorMessage = 'We couldn’t sign you in. Please try again.';
			}
		});
	});
</script>

<main class="min-h-screen w-full flex flex-col items-center justify-center bg-zinc-900 text-white">
	{#if status === 'loading'}
		<p class="text-lg text-gray-300">Signing you in…</p>
	{:else if status === 'error'}
		<p class="text-lg text-red-400 mb-4">{errorMessage}</p>
		<a href="/login" class="text-violet-400 hover:text-violet-300 underline">Back to login</a>
	{/if}
</main>
