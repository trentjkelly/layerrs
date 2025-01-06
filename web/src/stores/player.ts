import { writable } from 'svelte/store';

// Current audio information for all components to be able to read/write for the session
export const currentSong = writable(null);
export const isPlaying = writable(false);

// For the sidebar and bottom audio player to be visible
export const isSidebarOpen = writable(true);
// export const isSongSelected = writable(true);
