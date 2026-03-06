import { writable, get } from 'svelte/store';
import { jwt } from './auth';
import { urlBase } from './environment';
import { logger } from '../modules/lib/logger';

export const username = writable('');
export const email = writable('');
export const bio = writable('');
export const portraitUrl = writable('');

export async function loadProfile(): Promise<void> {
    if (get(username) !== '') return;

    try {
        const res = await fetch(`${get(urlBase)}/api/profile`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${get(jwt)}` }
        });
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
