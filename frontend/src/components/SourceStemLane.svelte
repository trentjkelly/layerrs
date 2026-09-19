<script lang="ts">
    import { onMount } from 'svelte';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase } from '../stores/environment';
    import type { TrackInfo } from '../models/types';
    import { audio, currentTrack, currentTrackId, isPlaying, currentTime as globalCurrentTime } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';
    import { trackPlayProgress } from '../modules/lib/play-tracking';
    import WaveformBar from './WaveformBar.svelte';
    import PlayPauseButton from './PlayPauseButton.svelte';

    let { track }: { track: TrackInfo } = $props();

    let urlBase = $state('');
    let newAudioURL = $state('');
    let waveformWidth = $state(0);
    let visibleBars = $state<number[]>([]);
    let timePerBar = $state(0);
    let currentTime = $state(0);
    let cursorTime = $state(0);
    let isCursorHovered = $state(false);
    let currentSongPercentage = $state(0);

    $effect(() => {
        updateWaveformBars();
    });

    $effect(() => {
        if (!$audio) return;

        const updateTime = () => {
            if ($currentTrackId === track.id) {
                currentTime = $audio.currentTime;
                globalCurrentTime.set(currentTime);
                currentSongPercentage = currentTime / track.duration * 100;
            }
        };

        const handleEnded = () => {
            if ($currentTrackId === track.id) {
                currentTime = 0;
                globalCurrentTime.set(0);
                currentSongPercentage = 0;
            }
        };

        $audio.addEventListener('timeupdate', updateTime);
        $audio.addEventListener('ended', handleEnded);
        return () => {
            $audio.removeEventListener('timeupdate', updateTime);
            $audio.removeEventListener('ended', handleEnded);
        };
    });

    $effect(() => {
        if ($currentTrackId === track.id && urlBase) {
            return trackPlayProgress(track.id, urlBase);
        }
    });

    onMount(() => {
        urlBase = getUrlBase();
    });

    function updateWaveformBars() {
        const barStride = 4;
        const numBars = Math.floor(waveformWidth / barStride);
        if (numBars <= 0 || track.waveformData.length === 0) {
            visibleBars = [];
            timePerBar = 0;
            return;
        }

        visibleBars = Array.from({ length: numBars }, (_, index) => {
            const waveformIndex = Math.min(track.waveformData.length - 1, Math.floor(track.waveformData.length * index / numBars));
            return track.waveformData[waveformIndex];
        });
        timePerBar = track.duration / numBars;
    }

    function handleMouseMove(event: MouseEvent) {
        const bounds = (event.currentTarget as HTMLElement).getBoundingClientRect();
        const percentage = Math.max(0, Math.min(1, (event.clientX - bounds.left) / waveformWidth));
        cursorTime = percentage * track.duration;
    }

    async function playOrSeek() {
        if (!$audio) return;

        if ($currentTrackId === track.id) {
            isPlaying.set(true);
            $audio.currentTime = cursorTime;
            await $audio.play();
            return;
        }

        isPlaying.set(true);
        newAudioURL = await getAudio(urlBase, String(track.id));
        currentTrack.set(newAudioURL);
        currentTrackId.set(track.id);
        $audio.src = newAudioURL;
        try {
            await $audio.play();
        } catch (error) {
            logger.error(`Failed to play source stem: ${error}`);
        }
    }

    function formatDuration(seconds: number): string {
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    }
</script>

<div class="relative flex min-w-0 items-stretch bg-zinc-950/50">
    <div class="absolute -top-3 bottom-1/2 left-0 w-px bg-primary"></div>
    <div class="absolute left-0 top-1/2 h-px w-4 bg-primary"></div>
    <div class="flex min-w-0 items-center gap-2 py-2 pl-6 pr-3 sm:w-64 sm:shrink-0 sm:border-r sm:border-zinc-800">
        <span class="text-[9px] uppercase tracking-wide text-primary-text">Stem</span>
        <img src={track.artistPortraitUrl} alt="{track.artistName} profile" class="h-7 w-7 shrink-0 object-cover">
        <div class="min-w-0 flex-1">
            <a href="/track/{track.id}" class="block truncate text-xs text-zinc-200 hover:text-white hover:underline">{track.description}</a>
            <a href="/profile/{track.artistId}" class="block truncate text-[10px] text-zinc-500 hover:text-white hover:underline">{track.artistName}</a>
        </div>
        <span class="shrink-0 text-[10px] text-zinc-500">{formatDuration(track.duration)}</span>
    </div>
    <div class="flex shrink-0 items-center border-r border-zinc-800 px-2">
        <PlayPauseButton {track} />
    </div>
    <div
        class="relative flex h-10 min-w-0 flex-1 cursor-pointer items-center overflow-hidden bg-primary bg-[repeating-linear-gradient(90deg,rgba(0,0,0,0.16)_0,rgba(0,0,0,0.16)_1px,transparent_1px,transparent_48px)] px-2"
        bind:clientWidth={waveformWidth}
        onmousemove={handleMouseMove}
        onmouseenter={() => isCursorHovered = true}
        onmouseleave={() => isCursorHovered = false}
        onclick={playOrSeek}
        onkeydown={(event) => event.key === 'Enter' && playOrSeek()}
        role="button"
        tabindex="0"
        aria-label="Play or seek through source stem {track.description}"
    >
        <div class="pointer-events-none absolute inset-y-0 w-px bg-black/40" style="left: {currentSongPercentage}%"></div>
        {#each visibleBars as bar, index}
            <WaveformBar
                height={bar}
                timePerBar={timePerBar}
                {index}
                cursorTime={cursorTime}
                {isCursorHovered}
                playedPercentage={currentSongPercentage / 100}
                trackDuration={track.duration}
            />
        {/each}
    </div>
    <div class="hidden shrink-0 border-l border-zinc-800 sm:block sm:w-44"></div>
</div>
