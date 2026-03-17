<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { urlBase } from '../../../stores/environment';
	import TrackNodeCard from '../../../components/TrackNodeCard.svelte';
	import { isSidebarOpen } from '../../../stores/player';
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
	};

	type SimNode = TrackInfo & { x?: number; y?: number; vx?: number; vy?: number; fx?: number | null; fy?: number | null; index?: number };

	type NodePos = { id: number; x: number; y: number };
	type LinkPos = { x1: number; y1: number; x2: number; y2: number };

	const CARD_W = 192;  // w-48
	const CARD_H = 96;   // h-24
	const LINK_OFFSET = 55;       // CARD_H/2 + small gap
	const SVG_HEIGHT = 1200;
	const LEVEL_SPACING = 200;    // CARD_H + ~100px gap

	let tracks: TrackInfo[] = [];
	let trackGraph: TrackTree[] | null = null;
	let idSet = new Set<number>();
	let trackMap = $state(new Map<number, TrackInfo>());

	let containerWidth = $state(800);
	let nodePositions = $state<NodePos[]>([]);
	let linkPositions = $state<LinkPos[]>([]);

	onMount(async () => {
		const slug = page.params.slug;
		const trackId = Number(slug);
		await fetchGraph(trackId);
		await fetchTracks([...idSet]);
		if (tracks.length > 0) initGraph();
	});

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

<main class={`transition-all duration-300 h-screen w-full overflow-y-auto ${$isSidebarOpen ? 'ml-64' : 'ml-0'} bg-zinc-900`}>
	<TopHeader pageName="Track" pageIcon="" />
	<section class="w-full" bind:clientWidth={containerWidth}>
		<svg width={containerWidth} height={SVG_HEIGHT}>
			<defs>
				<marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
					<path d="M 0 0 L 10 5 L 0 10 z" fill="#52525b" />
				</marker>
			</defs>

			<!-- Edges -->
			{#each linkPositions as link}
				<line
					x1={link.x1} y1={link.y1}
					x2={link.x2} y2={link.y2}
					stroke="#3f3f46"
					stroke-width="1.5"
					marker-end="url(#arrow)"
				/>
			{/each}

			<!-- Nodes -->
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
</main>
