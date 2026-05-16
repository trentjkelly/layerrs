<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { urlBase } from '../../../stores/environment';
	import TrackNodeCard from '../../../components/TrackNodeCard.svelte';
	import { isSidebarOpen } from '../../../stores/player';
	import { isLoggedIn } from '../../../stores/auth';
	import { fetchWithAuth } from '../../../modules/lib/fetch';
	import TopHeader from '../../../components/TopHeader.svelte';
	import * as d3 from 'd3';

	type TrackTree = {
		rootId: number;
		childId: number;
		derivationTag: 'layerr' | 'stem';
	};

	type TrackInfo = {
		id: number;
		artistId: number;
		likes: number;
		layerrs: number;
		description: string;
		artistName: string;
		artistPortraitUrl: string;
		color: string;
		duration: number;
		waveformData: number[];
		isLiked: boolean;
		priceInCents: number;
	};

	type SimNode = TrackInfo & { x?: number; y?: number; vx?: number; vy?: number; fx?: number | null; fy?: number | null; index?: number };

	type NodePos = { id: number; x: number; y: number };
	type LinkPos = { x1: number; y1: number; x2: number; y2: number };

	const CARD_W = 192;
	const CARD_H = 96;
	const LINK_OFFSET = 55;
	const SVG_HEIGHT = 1200;
	const LEVEL_SPACING = 200;

	let tracks: TrackInfo[] = [];
	let trackGraph: TrackTree[] | null = null;
	let idSet = new Set<number>();
	let trackMap = $state(new Map<number, TrackInfo>());

	let rootTrack = $state<TrackInfo | null>(null);
	let containerWidth = $state(800);
	let nodePositions = $state<NodePos[]>([]);
	let linkPositions = $state<LinkPos[]>([]);

	let showCheckoutModal = $state(false);
	let isProcessingPayment = $state(false);
	let paymentError = $state<string | null>(null);
	let downloadUrl = $state<string | null>(null);
	let downloadExpirationMinutes = $state<number>(0);

	onMount(async () => {
		const slug = page.params.slug;
		const trackId = Number(slug);
		await fetchGraph(trackId);
		await fetchTracks([...idSet]);
		if (tracks.length > 0) {
			rootTrack = tracks.find(t => t.id === trackId) ?? tracks[0];
			initGraph();
		}
	});

	async function startCheckout() {
		if (!$isLoggedIn) {
			window.location.href = '/login';
			return;
		}

		if (!rootTrack || rootTrack.priceInCents <= 0) return;

		showCheckoutModal = true;
		isProcessingPayment = true;
		paymentError = null;
		downloadUrl = null;

		try {
			const checkoutRes = await fetchWithAuth(`${$urlBase}/api/purchase/checkout`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ trackId: rootTrack.id })
			});

			const checkoutData = await checkoutRes.json();
			if (!checkoutRes.ok) {
				throw new Error(checkoutData.error || 'Failed to start checkout');
			}

			const stripe = (window as any).Stripe;
			if (!stripe) {
				throw new Error('Stripe not loaded');
			}

			const { error, paymentIntent } = await stripe.confirmCardPayment(checkoutData.clientSecret, {
				payment_method: 'pm_card_visa'
			});

			if (error) {
				throw new Error(error.message);
			}

			if (paymentIntent.status === 'succeeded') {
				const confirmRes = await fetchWithAuth(`${$urlBase}/api/purchase/confirm`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ paymentIntentId: paymentIntent.id })
				});

				const confirmData = await confirmRes.json();
				if (!confirmRes.ok) {
					throw new Error(confirmData.error || 'Failed to confirm purchase');
				}

				downloadUrl = confirmData.downloadURL;
				downloadExpirationMinutes = confirmData.downloadUrlExpirationMinutes;
			}
		} catch (err) {
			paymentError = err instanceof Error ? err.message : 'An error occurred';
		} finally {
			isProcessingPayment = false;
		}
	}

	function formatPrice(cents: number): string {
		return (cents / 100).toFixed(2);
	}

	function computeDepths(nodes: SimNode[]): Map<number, number> {
		const parents = new Map<number, Set<number>>();
		const children = new Map<number, Set<number>>();
		for (const n of nodes) { parents.set(n.id, new Set()); children.set(n.id, new Set()); }
		for (const e of trackGraph ?? []) {
			parents.get(e.childId)?.add(e.rootId);
			children.get(e.rootId)?.add(e.childId);
		}
		const depths = new Map<number, number>();
		const queue: number[] = [];
		for (const n of nodes) {
			if (parents.get(n.id)!.size === 0) { depths.set(n.id, 0); queue.push(n.id); }
		}
		while (queue.length > 0) {
			const id = queue.shift()!;
			const depth = depths.get(id)!;
			for (const childId of children.get(id) ?? []) {
				if ((depths.get(childId) ?? -1) < depth + 1) {
					depths.set(childId, depth + 1);
					queue.push(childId);
				}
			}
		}
		return depths;
	}

	function initGraph() {
		const nodes: SimNode[] = tracks.map(t => ({ ...t }));
		const nodeById = new Map(nodes.map(n => [n.id, n]));

		const links = (trackGraph ?? [])
			.filter(e => nodeById.has(e.rootId) && nodeById.has(e.childId))
			.map(e => ({
				source: nodeById.get(e.rootId)!,
				target: nodeById.get(e.childId)!,
				derivationTag: e.derivationTag,
			}));

		const depths = computeDepths(nodes);
		for (const node of nodes) {
			node.fy = (depths.get(node.id) ?? 0) * LEVEL_SPACING + LEVEL_SPACING;
		}

		d3.forceSimulation(nodes)
			.force('link', d3.forceLink(links).id((d: any) => d.id).distance(400))
			.force('charge', d3.forceManyBody().strength(-800))
			.force('x', d3.forceX(containerWidth / 2).strength(0.08))
			.force('collision', d3.forceCollide(110))
			.on('tick', () => {
				nodePositions = nodes.map(n => ({ id: n.id, x: n.x ?? 0, y: n.y ?? 0 }));
				linkPositions = links.map(l => {
					const s = l.source as SimNode;
					const t = l.target as SimNode;
					return shortenLink(s.x ?? 0, s.y ?? 0, t.x ?? 0, t.y ?? 0, LINK_OFFSET);
				});
			});
	}

	function shortenLink(x1: number, y1: number, x2: number, y2: number, offset: number): LinkPos {
		const dx = x2 - x1, dy = y2 - y1;
		const len = Math.sqrt(dx * dx + dy * dy);
		if (len < offset * 2) return { x1, y1, x2, y2 };
		const ux = dx / len, uy = dy / len;
		return { x1: x1 + ux * offset, y1: y1 + uy * offset, x2: x2 - ux * offset, y2: y2 - uy * offset };
	}

	async function fetchGraph(trackId: number) {
		const res = await fetch(`${$urlBase}/api/track/${trackId}/graph`);
		if (!res.ok) {
			console.error('Failed to fetch track graph');
			return;
		}
		trackGraph = (await res.json() as TrackTree[]).map(t => ({
			rootId: Number(t.rootId),
			childId: Number(t.childId),
			derivationTag: t.derivationTag,
		}));
		idSet = new Set<number>([trackId]);
		trackGraph.forEach(t => {
			idSet.add(t.rootId);
			idSet.add(t.childId);
		});
	}

	async function fetchTracks(trackIds: number[]) {
		const res = await fetch(`${$urlBase}/api/track/batch`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ track_ids: trackIds }),
		});
		if (!res.ok) {
			console.error('Failed to fetch track batch');
			return;
		}
		tracks = await res.json();
		trackMap = new Map(tracks.map(t => [t.id, t]));
	}
 </script>

<main class={`transition-all duration-300 h-screen w-full overflow-y-auto ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-olive-400`}>
	<TopHeader pageName="Track" pageIcon="" />

	{#if rootTrack && rootTrack.priceInCents > 0}
		<div class="fixed bottom-8 left-1/2 -translate-x-1/2 z-50">
			<button
				onclick={startCheckout}
				class="px-8 py-4 bg-emerald-600 hover:bg-emerald-700 rounded-full text-white font-semibold text-lg transition-colors shadow-lg flex items-center gap-3"
			>
				<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
				</svg>
				Buy for ${formatPrice(rootTrack.priceInCents)}
			</button>
		</div>
	{/if}

	<section class="w-full" bind:clientWidth={containerWidth}>
		<svg width={containerWidth} height={SVG_HEIGHT}>
			<defs>
				<marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
					<path d="M 0 0 L 10 5 L 0 10 z" fill="#52525b" />
				</marker>
			</defs>

			{#each linkPositions as link}
				<line
					x1={link.x1} y1={link.y1}
					x2={link.x2} y2={link.y2}
					stroke="#3f3f46"
					stroke-width="1.5"
					marker-end="url(#arrow)"
				/>
			{/each}

			{#each nodePositions as pos}
				<foreignObject
					x={pos.x - CARD_W / 2}
					y={pos.y - CARD_H / 2}
					width={CARD_W}
					height={CARD_H}
				>
					<div>
						{#if trackMap.has(pos.id)}
							<TrackNodeCard track={trackMap.get(pos.id)!} />
						{/if}
					</div>
				</foreignObject>
			{/each}
		</svg>
	</section>

	{#if showCheckoutModal}
		<div class="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
			<div class="bg-zinc-800 rounded-2xl p-8 max-w-md w-full">
				<h2 class="text-2xl font-bold text-white mb-4">
					{#if downloadUrl}
						Purchase Complete!
					{:else}
						Purchase Track
					{/if}
				</h2>

				{#if downloadUrl}
					<div class="space-y-4">
						<p class="text-zinc-300">Your download is ready. The link will expire in {downloadExpirationMinutes} minutes.</p>
						<a
							href={downloadUrl}
							download
							class="block w-full py-3 bg-emerald-600 hover:bg-emerald-700 rounded-lg text-white font-semibold text-center transition-colors"
						>
							Download WAV File
						</a>
						<button
							onclick={() => { showCheckoutModal = false; downloadUrl = null; }}
							class="w-full py-2 text-zinc-400 hover:text-white transition-colors"
						>
							Close
						</button>
					</div>
				{:else if isProcessingPayment}
					<div class="flex items-center justify-center py-8">
						<svg class="animate-spin h-8 w-8 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
						</svg>
					</div>
				{:else if paymentError}
					<div class="space-y-4">
						<p class="text-red-400">{paymentError}</p>
						<button
							onclick={startCheckout}
							class="w-full py-3 bg-indigo-600 hover:bg-indigo-700 rounded-lg text-white font-semibold transition-colors"
						>
							Try Again
						</button>
						<button
							onclick={() => { showCheckoutModal = false; paymentError = null; }}
							class="w-full py-2 text-zinc-400 hover:text-white transition-colors"
						>
							Cancel
						</button>
					</div>
				{:else}
					<p class="text-zinc-300 mb-4">Processing payment...</p>
				{/if}
			</div>
		</div>
	{/if}
</main>
