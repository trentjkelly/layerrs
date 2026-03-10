import { writable, get } from 'svelte/store';
import { urlBase } from './environment';
import { logger } from '../modules/lib/logger';
import { fetchWithAuth } from '../modules/lib/fetch';

export const username = writable('');
export const email = writable('');
export const bio = writable('');
export const portraitUrl = writable('');

export async function loadProfile(): Promise<void> {
    if (get(username) !== '') return;

    try {
        const res = await fetchWithAuth(`${get(urlBase)}/api/profile`, { method: 'GET' });
        if (res.ok) {
            const data = await res.json();
            username.set(data.username ?? '');
            email.set(data.email ?? '');
            bio.set(data.bio ?? '');
            portraitUrl.set(data.portraitUrl ?? '');
        }
    } catch (err) {
        logger.error(`Failed to load profile: ${err}`);
    }
}
