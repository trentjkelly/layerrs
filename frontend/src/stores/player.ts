import { writable } from 'svelte/store';

// Current audio information for all components to be able to read/write for the session
export const currentTrack = writable<string | null>('');
export const isPlaying = writable(false);
export const currentTrackId = writable(0);

// For the sidebar and bottom audio player to be visible
export const isSidebarOpen = writable(true);
// export const isSongSelected = writable(true);
export const audio = writable<HTMLAudioElement | null>(null);

