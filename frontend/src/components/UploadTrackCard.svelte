<script lang="ts">
    import { audio, currentTrackId, isPlaying } from '../stores/player';
    import { urlBase } from '../stores/environment';
    import { getAudio } from '../modules/requests/track-requests';

    type TrackPreview = {
        id: number;
        description: string;
        artistName: string;
        artistPortraitUrl: string;
    };

    let { track }: { track: TrackPreview } = $props();

    async function playTrack() {
        if (!$audio) return;

        if ($currentTrackId === track.id) {
            if ($isPlaying) {
                $audio.pause();
                isPlaying.set(false);
            } else {
                await $audio.play();
                isPlaying.set(true);
            }
            return;
        }

        const url = await getAudio($urlBase, String(track.id));
        if (!url) return;

        currentTrackId.set(track.id);
        $audio.src = url;
        await $audio.play();
        isPlaying.set(true);
    }
</script>

<div class="flex items-center gap-3 px-4 py-3 text-white">
    <button
        type="button"
        onclick={playTrack}
        class="shrink-0 w-7 h-7 flex items-center justify-center text-white hover:text-violet-400 transition-colors"
    >
        {#if $currentTrackId === track.id && $isPlaying}
            <!-- Pause icon -->
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-5 h-5">
                <path fill-rule="evenodd" d="M6.75 5.25a.75.75 0 0 1 .75-.75H9a.75.75 0 0 1 .75.75v13.5a.75.75 0 0 1-.75.75H7.5a.75.75 0 0 1-.75-.75V5.25zm7.5 0A.75.75 0 0 1 15 4.5h1.5a.75.75 0 0 1 .75.75v13.5a.75.75 0 0 1-.75.75H15a.75.75 0 0 1-.75-.75V5.25z" clip-rule="evenodd" />
            </svg>
        {:else}
            <!-- Play icon -->
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-5 h-5">
                <path fill-rule="evenodd" d="M4.5 5.653c0-1.427 1.529-2.33 2.779-1.643l11.54 6.347c1.295.712 1.295 2.573 0 3.286L7.28 19.99c-1.25.687-2.779-.217-2.779-1.643V5.653z" clip-rule="evenodd" />
            </svg>
        {/if}
    </button>
    {#if track.artistPortraitUrl}
        <img src={track.artistPortraitUrl} alt={track.artistName} class="w-8 h-8 rounded-xl object-cover shrink-0" />
    {:else}
        <div class="w-8 h-8 rounded-xl bg-zinc-600 shrink-0"></div>
    {/if}
    <div class="flex items-center gap-2 min-w-0">
        <p class="text-base font-medium text-zinc-400 truncate">{track.artistName}</p>
        <span class="text-zinc-500">·</span>
        <p class="text-base truncate">"{track.description}"</p>
    </div>
</div>
