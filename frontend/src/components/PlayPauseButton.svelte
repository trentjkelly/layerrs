<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { logger } from '../modules/lib/logger';
    import { getUrlBase } from '../stores/environment';
    import type { TrackInfo } from '../models/types';
    import { audio, currentTrack, currentTrackId, isPlaying } from '../stores/player';
    import { getAudio } from '../modules/requests/track-requests';

    let {
        track
    }: {
        track: TrackInfo;
    } = $props();

    let newAudioURL = $state('');
    let urlBase = $state('');

    let isCurrentTrackPlaying = $derived($currentTrackId === track.id && $isPlaying);

    onMount(async () => {
        urlBase = getUrlBase();
    });

    onDestroy(() => {
        if (newAudioURL) {
            URL.revokeObjectURL(newAudioURL);
        }
    });

    async function playPauseAudio() {
        if ($audio) {
            if (track.id === $currentTrackId) {
                if ($audio.paused) {
                    isPlaying.set(true);
                    await $audio.play();
                } else {
                    isPlaying.set(false);
                    $audio.pause();
                }
            } else {
                if (!$audio.paused) {
                    isPlaying.set(false);
                    $audio.pause();
                }

                isPlaying.set(true);
                newAudioURL = await getAudio(urlBase, String(track.id));
                currentTrack.set(newAudioURL);
                currentTrackId.set(track.id);

                $audio.src = newAudioURL;
                try {
                    await $audio.play();
                } catch (error) {
                    logger.error(`Failed to play audio: ${error}`);
                }
            }
        }
    }
</script>

<button
    class="flex items-center justify-center hover:cursor-pointer shrink-0 p-2 bg-olive-700 hover:bg-olive-800 transition-colors"
    onclick={playPauseAudio}
    aria-label={isCurrentTrackPlaying ? 'Pause' : 'Play'}
>
    {#if isCurrentTrackPlaying}
        <img src="pause.png" alt="Pause" class="w-3.5 h-3.5">
    {:else}
        <img src="play.png" alt="Play" class="w-3.5 h-3.5">
    {/if}
</button>
