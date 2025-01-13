import { writable } from "svelte/store";

// JSON web token, written on page load from cookies
export const jwt = writable('');
export const refreshToken = writable('');