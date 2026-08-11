<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { handleBrowserLogin } from '../../../modules/lib/session';
	import { loadProfile } from '../../../stores/profile';
	import { logger } from '../../../modules/lib/logger';

	let status = $state<'loading' | 'success' | 'error'>('loading');
	let errorMessage = $state('');

	onMount(() => {
		const hash = typeof window !== 'undefined' ? window.location.hash.slice(1) : '';
		const params = new URLSearchParams(hash);
		const jwt = params.get('jwt');
		const refresh = params.get('refresh');
		const firstLogin = params.get('firstLogin') === 'true';

		if (!jwt || !refresh) {
			logger.error('Callback missing jwt or refresh in URL hash');
			status = 'error';
			errorMessage = 'Invalid sign-in link. Please request a new one from the login page.';
			return;
		}

		handleBrowserLogin(jwt, refresh).then(async (success) => {
			if (success) {
				status = 'success';
				await loadProfile();
				goto(firstLogin ? '/profile' : '/');
			} else {
				status = 'error';
				errorMessage = 'We couldn’t sign you in. Please try again.';
			}
		});
	});
</script>

<main class="min-h-screen w-full flex flex-col items-center justify-center bg-zinc-950 text-white">
	{#if status === 'loading'}
		<p class="text-lg text-zinc-300">Signing you in…</p>
	{:else if status === 'error'}
		<p class="text-lg text-danger mb-4">{errorMessage}</p>
		<a href="/login" class="text-primary-text hover:text-primary-hover underline">Back to login</a>
	{/if}
</main>
